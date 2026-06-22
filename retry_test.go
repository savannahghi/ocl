package ocl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newRetryClient builds a Client pointed at srv with the given policy so the
// real WithRetry wiring (and the retryRoundTripper it installs) is exercised.
func newRetryClient(t *testing.T, baseURL string, policy RetryPolicy) *Client {
	t.Helper()

	c, err := NewClient(baseURL, "test-token", WithRetry(policy))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	return c
}

func TestWithRetry(t *testing.T) {
	constStatus := func(code int) func(w http.ResponseWriter, callNum int32) {
		return func(w http.ResponseWriter, _ int32) {
			w.WriteHeader(code)
		}
	}

	tests := []struct {
		name       string
		policy     RetryPolicy
		method     string
		timeout    time.Duration
		handler    func(w http.ResponseWriter, callNum int32)
		wantErr    bool
		wantStatus int
		checkCalls func(t *testing.T, calls int32)
	}{
		{
			name:   "retries transient 5xx then succeeds",
			policy: RetryPolicy{MaxAttempts: 3, InitialBackoff: time.Millisecond},
			method: http.MethodGet,
			handler: func(w http.ResponseWriter, callNum int32) {
				if callNum < 3 {
					w.WriteHeader(http.StatusServiceUnavailable)

					return
				}

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"ok":true}`))
			},
			wantStatus: http.StatusOK,
			checkCalls: func(t *testing.T, calls int32) {
				t.Helper()

				if calls != 3 {
					t.Fatalf("want 3 attempts, got %d", calls)
				}
			},
		},
		{
			name:       "does not retry non-retryable status",
			policy:     RetryPolicy{MaxAttempts: 3, InitialBackoff: time.Millisecond},
			method:     http.MethodGet,
			handler:    constStatus(http.StatusBadRequest),
			wantStatus: http.StatusBadRequest,
			checkCalls: func(t *testing.T, calls int32) {
				t.Helper()

				if calls != 1 {
					t.Fatalf("400 must not retry: want 1 call, got %d", calls)
				}
			},
		},
		{
			name:       "does not retry non-idempotent POST",
			policy:     RetryPolicy{MaxAttempts: 3, InitialBackoff: time.Millisecond},
			method:     http.MethodPost,
			handler:    constStatus(http.StatusServiceUnavailable),
			wantStatus: http.StatusServiceUnavailable,
			checkCalls: func(t *testing.T, calls int32) {
				t.Helper()

				if calls != 1 {
					t.Fatalf("POST must not retry: want 1 call, got %d", calls)
				}
			},
		},
		{
			name:    "stops on context deadline",
			policy:  RetryPolicy{MaxAttempts: 5, InitialBackoff: 200 * time.Millisecond},
			method:  http.MethodGet,
			timeout: 20 * time.Millisecond,
			handler: constStatus(http.StatusServiceUnavailable),
			wantErr: true,
			checkCalls: func(t *testing.T, calls int32) {
				t.Helper()

				if calls >= 5 {
					t.Fatalf("context deadline must cut retries short: got %d calls", calls)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls int32

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				tt.handler(w, atomic.AddInt32(&calls, 1))
			}))
			defer srv.Close()

			c := newRetryClient(t, srv.URL, tt.policy)

			ctx := context.Background()
			if tt.timeout > 0 {
				var cancel context.CancelFunc

				ctx, cancel = context.WithTimeout(ctx, tt.timeout)
				defer cancel()
			}

			req, _ := http.NewRequestWithContext(ctx, tt.method, srv.URL, http.NoBody)

			resp, err := c.HTTP.Do(req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				defer resp.Body.Close()

				if tt.wantStatus != 0 && resp.StatusCode != tt.wantStatus {
					t.Fatalf("want status %d, got %d", tt.wantStatus, resp.StatusCode)
				}
			}

			tt.checkCalls(t, atomic.LoadInt32(&calls))
		})
	}
}
