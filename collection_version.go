package ocl

import (
	"context"
	"net/http"
)

func (c *Client) CreateCollectionVersion(
	ctx context.Context, collectionVersion *CollectionVersion, headers *Headers,
) (*CollectionVersion, error) {
	var resp *CollectionVersion

	err := c.makeRequest(ctx, http.MethodPost, composeCollectionVersionsURL(headers), nil, collectionVersion, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// composeCollectionVersionsURL forms the create/get collection versions url. It follows this path
// /orgs/:org/collections/:collection/versions/.
func composeCollectionVersionsURL(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/collections/" + headers.Collection + "/versions/"
}
