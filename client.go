package ocl

import (
	"errors"
	"net"
	"net/http"
	"os"
	"time"
)

// DefaultHTTPTimeout caps the wall-clock time of any single OCL request.
// It is intentionally tight because this SDK is typically used inside
// services serving live user traffic; failing fast under upstream stress
// is preferable to holding connections open for long periods.
const DefaultHTTPTimeout = 10 * time.Second

// NewDefaultTransport returns an http.Transport tuned for server-side use
// of this SDK. The defaults that http.DefaultTransport ships with target
// CLI-style usage (MaxIdleConnsPerHost = 2), which causes connection-pool
// starvation in services that proxy high calls to OCL. Callers who
// need to customize the transport (e.g. wrap it with an OTel RoundTripper)
// can call this and mutate the returned value before passing it to
// WithHTTPClient.
func NewDefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 64,
		MaxConnsPerHost:     128,
		IdleConnTimeout:     90 * time.Second,

		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}
}

// NewDefaultHTTPClient returns an *http.Client built on NewDefaultTransport
// with DefaultHTTPTimeout. Use this as a starting point if you want to
// customize the client before passing it to WithHTTPClient.
func NewDefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: NewDefaultTransport(),
		Timeout:   DefaultHTTPTimeout,
	}
}

type Client struct {
	baseURL string
	token   string

	HTTP *http.Client
}

// ClientOption allows customization of the client.
type ClientOption func(c *Client)

// WithTimeout sets the total per-request wall-clock timeout enforced by
// the underlying http.Client. Callers using context-based deadlines should
// still keep this set as a defense-in-depth ceiling.
func WithTimeout(t time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTP.Timeout = t
	}
}

// WithHTTPClient replaces the default HTTP client. Use this when you need
// to inject a custom transport (e.g. for tracing, mTLS, or testing). If
// you only need to tweak the default transport, prefer:
//
//	t := ocl.NewDefaultTransport()
//	t.MaxIdleConnsPerHost = 128
//	c, _ := ocl.NewClient(url, token, ocl.WithHTTPClient(&http.Client{
//	    Transport: t,
//	    Timeout:   10 * time.Second,
//	}))
//
// A nil client is ignored so callers don't accidentally disable the SDK
// by passing through an unset value.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.HTTP = httpClient
		}
	}
}

// NewClientFromEnvVars creates a new client where the needed fields are
// retrieved from the environment variables.
func NewClientFromEnvVars() (*Client, error) {
	return NewClient(os.Getenv("OCL_BASE_URL"), os.Getenv("OCL_TOKEN"))
}

// NewClient creates a new ocl api client. The returned client uses a
// server-tuned HTTP transport by default (see NewDefaultTransport); use
// WithHTTPClient to override.
func NewClient(baseURL string, token string, options ...ClientOption) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is empty")
	}

	if token == "" {
		return nil, errors.New("token is empty")
	}

	client := &Client{
		HTTP:    NewDefaultHTTPClient(),
		baseURL: baseURL,
		token:   token,
	}

	for _, opt := range options {
		opt(client)
	}

	return client, nil
}
