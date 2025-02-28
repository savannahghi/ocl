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

type CollectionReferenceResponse []struct {
	Message    string `json:"message"`
	Added      bool   `json:"added"`
	Expression string `json:"expression"`
}

func (c *Client) CreateCollectionReference(
	ctx context.Context, collectionRef *CollectionReference, headers *Headers,
) (*CollectionReferenceResponse, error) {
	var resp *CollectionReferenceResponse

	err := c.makeRequest(ctx, http.MethodPut, composeCollectionReferencesURL(headers), nil, collectionRef, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// composeCollectionVersionsURL forms the create/get collection references url. It follows this path
// PUT /orgs/:org/collections/:collection/references/.
func composeCollectionReferencesURL(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/collections/" + headers.Collection + "/references/"
}
