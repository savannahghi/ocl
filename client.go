package ocl

import (
	"errors"
	"net/http"
	"os"
	"time"
)

type Client struct {
	baseURL string
	token   string

	HTTP *http.Client
}

// ClientOption allows customization of the client.
type ClientOption func(c *Client)

func WithTimeout(t time.Duration) func(c *Client) {
	return func(c *Client) {
		c.HTTP.Timeout = t
	}
}

// NewClientFromEnvVars creates a new client where the needed fields are
// retrieved from the environment variables.
func NewClientFromEnvVars() (*Client, error) {
	return NewClient(os.Getenv("OCL_BASE_URL"), os.Getenv("OCL_TOKEN"))
}

// NewClient creates a new ocl api client.
func NewClient(baseURL string, token string, options ...ClientOption) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is empty")
	}

	if token == "" {
		return nil, errors.New("token is empty")
	}

	client := &Client{
		HTTP:    &http.Client{},
		baseURL: baseURL,
		token:   token,
	}

	for _, opt := range options {
		opt(client)
	}

	return client, nil
}
