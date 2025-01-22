package ocl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) newRequest(
	ctx context.Context,
	method, path string,
	params url.Values,
	data interface{},
) (*http.Request, error) {
	url, err := c.composeRequestURL(path, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	c.setHeaders(request)

	switch payload := data.(type) {
	case nil:
		request.Body = nil
	case io.ReadCloser:
		request.Body = payload
	case io.Reader:
		request.Body = io.NopCloser(payload)
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		request.Body = io.NopCloser(bytes.NewReader(b))
	}

	return request, nil
}

func (c *Client) setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+c.token)
}

func (c *Client) composeRequestURL(path string, params url.Values) (string, error) {
	u, err := url.Parse(c.baseURL + "/" + path)
	if err != nil {
		return "", errors.New("url parse: " + err.Error())
	}

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func (c *Client) readResponse(response *http.Response, result interface{}) error {
	if response.Body == nil {
		return errors.New("response body is nil")
	}
	defer response.Body.Close()

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBytes, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshall body: %w", err)
	}

	return nil
}

func (c *Client) makeRequest(
	ctx context.Context,
	method, path string,
	params url.Values,
	data, result interface{},
) error {
	request, err := c.newRequest(ctx, method, path, params, data)
	if err != nil {
		return err
	}

	resp, err := c.HTTP.Do(request)
	if err != nil {
		return err
	}

	return c.readResponse(resp, result)
}
