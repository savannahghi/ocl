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
	data any,
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

		reader := bytes.NewReader(b)
		request.Body = io.NopCloser(reader)

		// This code is here to thwart an error:
		//
		// "http2: Transport: cannot retry err [http2: Transport received Server's graceful shutdown GOAWAY]
		// after Request.Body was written; define Request.GetBody to avoid this error"
		//
		// This snippet defines how to recreate the request body if it needs to be resent (e.g., on retry(ies)).
		// Inside this function, we create a fresh bytes.NewReader(b) (so the read offset is reset).
		// Then we again wrap it in io.NopCloser, as required by http.Request.
		request.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(b)), nil
		}
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

func (c *Client) readResponse(response *http.Response, result any) error {
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
	//nolint:unparam
	params url.Values,
	data, result any,
) error {
	request, err := c.newRequest(ctx, method, path, params, data)
	if err != nil {
		return err
	}

	resp, err := c.HTTP.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		apiErr := APIErrorResponse{}

		err = json.Unmarshal(body, &apiErr)
		if err != nil {
			return fmt.Errorf("failed (HTTP %d): %s", resp.StatusCode, string(body))
		}

		return &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    string(body),
			APIError:   apiErr,
		}
	}

	if result != nil {
		err = c.readResponse(resp, result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) streamRawData(
	ctx context.Context,
	method, path string,
	params url.Values,
) (io.ReadCloser, error) {
	request, err := c.newRequest(ctx, method, path, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make new reqeust: %w", err)
	}

	resp, err := c.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	// See https://docs.openconceptlab.org/en/latest/oclapi/apireference/exportapi.html
	// for response codes returned by OCL when requested export data is not ready
	if resp.StatusCode < 299 && resp.StatusCode > 200 {
		return nil, errors.New("requested data is not ready")
	}

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response from body: %w", err)
		}

		return nil, fmt.Errorf("failed to process request: %v", string(body))
	}

	return resp.Body, nil
}

// Headers represents the custom headers sent to the client. In OCL, concepts are namespaced
// with Organisations and sources e.g You can have WHO as an org, and many sources within that
// org e.g ICD-10, ICD-11.
//
// The idea is that the client should send the source, org, collection, concept as headers so that the library will
// correctly create the API URL.
type Headers struct {
	Organisation string `json:"organisation,omitempty"`
	Source       string `json:"source,omitempty"`
	Collection   string `json:"collection,omitempty"`
	ConceptID    string `json:"concept,omitempty"`
	VersionID    string `json:"version_id,omitempty"`
	MappingID    string `json:"mapping,omitempty"`
}
