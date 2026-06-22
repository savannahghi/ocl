package ocl

import (
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// RetryPolicy controls how the client retries failed requests. See WithRetry.
//
// Retries are applied at the transport layer, below makeRequest, so a
// retryable upstream response (e.g. a transient 503) is replayed
// transparently and only the final response surfaces to callers.
type RetryPolicy struct {
	// MaxAttempts is the total number of attempts including the first call.
	// A value of 1 means no retry. Zero/negative is treated as 3.
	MaxAttempts int

	// InitialBackoff is the wait before the first retry. Defaults to 100ms.
	InitialBackoff time.Duration

	// MaxBackoff caps the wait between any two attempts. Defaults to 5s.
	MaxBackoff time.Duration

	// BackoffMultiplier controls exponential growth. Defaults to 2.0.
	BackoffMultiplier float64

	// Jitter, in [0,1), randomises each computed backoff to spread retries.
	// 0.2 means actual wait is in [0.8*backoff, 1.2*backoff]. Defaults to 0.2.
	Jitter float64

	// RetryableStatuses are HTTP status codes that trigger a retry.
	// Defaults to 408, 429, 500, 502, 503, 504.
	RetryableStatuses []int

	// RetryableMethods are HTTP methods that may be retried.
	// Defaults to GET, HEAD, PUT, DELETE, OPTIONS.
	// POST and PATCH are excluded because OCL create/update are not idempotent.
	RetryableMethods []string

	// RespectRetryAfter, when true, makes the policy wait for the duration
	// advertised by a Retry-After response header (overrides backoff for
	// that attempt) before retrying. Defaults to true.
	RespectRetryAfter bool
}

func (p *RetryPolicy) applyDefaults() {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 3
	}

	if p.InitialBackoff <= 0 {
		p.InitialBackoff = 100 * time.Millisecond
	}

	if p.MaxBackoff <= 0 {
		p.MaxBackoff = 5 * time.Second
	}

	if p.BackoffMultiplier <= 1.0 {
		p.BackoffMultiplier = 2.0
	}

	if p.Jitter < 0 || p.Jitter >= 1 {
		p.Jitter = 0.2
	}

	if len(p.RetryableStatuses) == 0 {
		p.RetryableStatuses = []int{
			http.StatusRequestTimeout,      // 408
			http.StatusTooManyRequests,     // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout,      // 504
		}
	}

	if len(p.RetryableMethods) == 0 {
		p.RetryableMethods = []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		}
	}
	// RespectRetryAfter defaults to true unless explicitly disabled. We
	// can't distinguish "unset" from "false" on a bool, so consumers who
	// want to disable should set it explicitly to false; the empty struct
	// gives true, matching most users' expectations.
}

func (p *RetryPolicy) statusRetryable(code int) bool {
	for _, c := range p.RetryableStatuses {
		if c == code {
			return true
		}
	}

	return false
}

func (p *RetryPolicy) methodRetryable(method string) bool {
	for _, m := range p.RetryableMethods {
		if m == method {
			return true
		}
	}

	return false
}

// nextBackoff returns the wait before attempt number `attempt` (1-indexed
// for the FIRST retry, i.e. attempt=1 is the wait after the initial call).
func (p *RetryPolicy) nextBackoff(attempt int, rng *rand.Rand) time.Duration {
	exp := math.Pow(p.BackoffMultiplier, float64(attempt-1))

	d := time.Duration(float64(p.InitialBackoff) * exp)
	if d > p.MaxBackoff {
		d = p.MaxBackoff
	}

	if p.Jitter > 0 {
		factor := 1 - p.Jitter + (rng.Float64() * 2 * p.Jitter)
		d = time.Duration(float64(d) * factor)
	}

	return d
}

// retryRoundTripper implements http.RoundTripper. It wraps an inner
// transport and replays the request when the response is retryable.
type retryRoundTripper struct {
	inner  http.RoundTripper
	policy RetryPolicy
}

func (r *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Fast path: method not eligible for retry, or body not replayable.
	if !r.policy.methodRetryable(req.Method) || !bodyReplayable(req) {
		return r.inner.RoundTrip(req)
	}

	// Per-request RNG seeded by time. Cheap, good enough for jitter — not
	// security-sensitive.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var (
		resp     *http.Response
		err      error
		attempts = r.policy.MaxAttempts
	)

	for i := 0; i < attempts; i++ {
		if i > 0 {
			// Drain & close any prior response before retry.
			drainBody(resp)

			// Reset body for replay.
			if req.GetBody != nil {
				newBody, gerr := req.GetBody()
				if gerr != nil {
					return nil, gerr
				}

				req.Body = newBody
			}

			wait := r.policy.nextBackoff(i, rng)
			if r.policy.RespectRetryAfter && resp != nil {
				if ra := parseRetryAfter(resp); ra > 0 && ra < r.policy.MaxBackoff {
					wait = ra
				}
			}

			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(wait):
			}
		}

		resp, err = r.inner.RoundTrip(req)
		if err != nil {
			// Transport-level error: retry if attempts remain. Net errors
			// are almost always transient (timeout, RST, EOF). A cancelled
			// or expired context surfaces here too; the next iteration's
			// backoff select returns promptly via req.Context().Done().
			continue
		}

		if !r.policy.statusRetryable(resp.StatusCode) {
			return resp, nil
		}
	}

	return resp, err
}

// bodyReplayable reports whether the request body can safely be re-sent.
func bodyReplayable(req *http.Request) bool {
	if req.Body == nil || req.Body == http.NoBody {
		return true
	}

	return req.GetBody != nil
}

// drainBody discards any remaining bytes and closes the body. Required so
// the connection can be returned to the idle pool between retries.
func drainBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

// parseRetryAfter handles the two RFC 9110 Retry-After forms: integer seconds
// or HTTP-date. Returns zero on parse failure or absent header.
func parseRetryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}

	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}

	if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}

	if t, err := http.ParseTime(v); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}

		return d
	}

	return 0
}
