package ocl

import (
	"context"
	"net/http"
)

type CollectionReference struct {
	Data struct {
		Expressions []string `json:"expressions"`
	} `json:"data"`
}

func (c *Client) CreateCollectionReference(
	ctx context.Context, collection *Collection, headers *Headers,
) (*Collection, error) {
	var resp Collection

	err := c.makeRequest(ctx, http.MethodPut, composeCollectionReferencesURL(headers), nil, collection, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// composeCollectionVersionsURL forms the create/get collection references url. It follows this path
// PUT /orgs/:org/collections/:collection/references/.
func composeCollectionReferencesURL(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/collections/" + headers.Collection + "/references/"
}
