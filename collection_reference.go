package ocl

import (
	"context"
	"fmt"
	"net/http"
)

type Expression struct {
	Expression []string `json:"expressions,omitempty"`
}

type CollectionReference struct {
	Data Expression `json:"data"`
}

type CollectionReferenceAsyncResponse []struct {
	ID       string `json:"id,omitempty"`
	State    string `json:"state,omitempty"`
	Name     string `json:"name,omitempty"`
	Queue    string `json:"queue,omitempty"`
	Username string `json:"username,omitempty"`
	Task     string `json:"task,omitempty"`
}

func (c *Client) CreateCollectionReference(
	ctx context.Context, collectionRef *CollectionReference, headers *Headers,
) (*CollectionReferenceAsyncResponse, error) {
	var resp *CollectionReferenceAsyncResponse

	path := fmt.Sprintf(
		"orgs/%s/collections/%s/references/?cascade=sourcetoconcepts&async=true",
		headers.Organisation, headers.Collection,
	)

	err := c.makeRequest(ctx, http.MethodPut, path, nil, collectionRef, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
