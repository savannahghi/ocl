package ocl

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	customClient := &http.Client{Timeout: 99 * time.Second}

	tests := []struct {
		name     string
		baseURL  string
		token    string
		options  []ClientOption
		wantErr  bool
		validate func(t *testing.T, c *Client)
	}{
		{
			name:    "missing baseURL returns error",
			baseURL: "",
			token:   "token",
			wantErr: true,
		},
		{
			name:    "missing token returns error",
			baseURL: "http://example.com",
			token:   "",
			wantErr: true,
		},
		{
			name:    "default uses tuned transport, not http.DefaultTransport",
			baseURL: "http://example.com",
			token:   "token",
			validate: func(t *testing.T, c *Client) {
				tr, ok := c.HTTP.Transport.(*http.Transport)
				if !ok {
					t.Fatalf("Transport is %T; want *http.Transport (the SDK should not fall back to http.DefaultTransport)", c.HTTP.Transport)
				}

				if tr.MaxIdleConnsPerHost < 16 {
					t.Errorf("MaxIdleConnsPerHost = %d; want >= 16 (http.DefaultTransport ships 2, which starves the pool in server use)", tr.MaxIdleConnsPerHost)
				}
			},
		},
		{
			name:    "WithTimeout sets the timeout and preserves the tuned transport",
			baseURL: "http://example.com",
			token:   "token",
			options: []ClientOption{WithTimeout(5 * time.Second)},
			validate: func(t *testing.T, c *Client) {
				if c.HTTP.Timeout != 5*time.Second {
					t.Errorf("Timeout = %v; want 5s", c.HTTP.Timeout)
				}

				tr, ok := c.HTTP.Transport.(*http.Transport)
				if !ok {
					t.Fatalf("Transport is %T after WithTimeout; want tuned *http.Transport", c.HTTP.Transport)
				}

				if tr.MaxIdleConnsPerHost < 16 {
					t.Errorf("WithTimeout should not regress the transport: MaxIdleConnsPerHost = %d", tr.MaxIdleConnsPerHost)
				}
			},
		},
		{
			name:    "WithHTTPClient replaces the underlying client",
			baseURL: "http://example.com",
			token:   "token",
			options: []ClientOption{WithHTTPClient(customClient)},
			validate: func(t *testing.T, c *Client) {
				if c.HTTP != customClient {
					t.Error("WithHTTPClient did not replace the HTTP client")
				}

				if c.HTTP.Timeout != 99*time.Second {
					t.Errorf("Timeout = %v; want 99s", c.HTTP.Timeout)
				}
			},
		},
		{
			name:    "WithHTTPClient(nil) is ignored so the SDK stays usable",
			baseURL: "http://example.com",
			token:   "token",
			options: []ClientOption{WithHTTPClient(nil)},
			validate: func(t *testing.T, c *Client) {
				if c.HTTP == nil {
					t.Fatal("HTTP became nil after WithHTTPClient(nil); the SDK would panic on the next request")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.baseURL, tt.token, tt.options...)

			if (err != nil) != tt.wantErr {
				t.Fatalf("NewClient err = %v; wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if tt.validate != nil {
				tt.validate(t, c)
			}
		})
	}
}

func TestNewDefaultTransport(t *testing.T) {
	tests := []struct {
		name   string
		assert func(t *testing.T, tr *http.Transport)
	}{
		{
			name: "MaxIdleConnsPerHost is server-sized, not http.DefaultTransport's 2",
			assert: func(t *testing.T, tr *http.Transport) {
				if tr.MaxIdleConnsPerHost < 16 {
					t.Errorf("MaxIdleConnsPerHost = %d; want >= 16", tr.MaxIdleConnsPerHost)
				}
			},
		},
		{
			name: "MaxConnsPerHost exceeds MaxIdleConnsPerHost to allow burst above steady-state",
			assert: func(t *testing.T, tr *http.Transport) {
				if tr.MaxConnsPerHost <= tr.MaxIdleConnsPerHost {
					t.Errorf("MaxConnsPerHost (%d) must exceed MaxIdleConnsPerHost (%d)",
						tr.MaxConnsPerHost, tr.MaxIdleConnsPerHost)
				}
			},
		},
		{
			name: "HTTP/2 is requested via ALPN",
			assert: func(t *testing.T, tr *http.Transport) {
				if !tr.ForceAttemptHTTP2 {
					t.Error("ForceAttemptHTTP2 should be true so we multiplex when the upstream supports it")
				}
			},
		},
		{
			name: "TLS handshake has a timeout",
			assert: func(t *testing.T, tr *http.Transport) {
				if tr.TLSHandshakeTimeout == 0 {
					t.Error("TLSHandshakeTimeout is unset; a hung handshake will tie up a goroutine indefinitely")
				}
			},
		},
		{
			name: "idle connections expire so they don't accumulate forever",
			assert: func(t *testing.T, tr *http.Transport) {
				if tr.IdleConnTimeout == 0 {
					t.Error("IdleConnTimeout is unset; idle keep-alives will linger forever and exhaust FDs")
				}
			},
		},
		{
			name: "two calls return distinct instances so callers can mutate safely",
			assert: func(t *testing.T, tr *http.Transport) {
				other := NewDefaultTransport()
				if tr == other {
					t.Fatal("NewDefaultTransport returned the same pointer twice; mutating one would affect the other")
				}

				tr.MaxIdleConnsPerHost = 1
				if other.MaxIdleConnsPerHost == 1 {
					t.Fatal("mutating one transport leaked into another instance")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewDefaultTransport()
			if tr == nil {
				t.Fatal("NewDefaultTransport returned nil")
			}

			tt.assert(t, tr)
		})
	}
}

func TestNewDefaultHTTPClient(t *testing.T) {
	tests := []struct {
		name   string
		assert func(t *testing.T, c *http.Client)
	}{
		{
			name: "Timeout is set so a hanging request can't run forever",
			assert: func(t *testing.T, c *http.Client) {
				if c.Timeout == 0 {
					t.Error("Timeout is zero; defense-in-depth requires a hard ceiling even when callers use ctx")
				}
			},
		},
		{
			name: "Transport is the tuned one, not http.DefaultTransport",
			assert: func(t *testing.T, c *http.Client) {
				tr, ok := c.Transport.(*http.Transport)
				if !ok {
					t.Fatalf("Transport is %T; want *http.Transport", c.Transport)
				}

				if tr.MaxIdleConnsPerHost < 16 {
					t.Errorf("MaxIdleConnsPerHost = %d; want the tuned value", tr.MaxIdleConnsPerHost)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewDefaultHTTPClient()
			if c == nil {
				t.Fatal("NewDefaultHTTPClient returned nil")
			}

			tt.assert(t, c)
		})
	}
}
