package ocl

import (
	"context"
	"net/http"
	"time"
)

type CollectionVersion struct {
	Type               string     `json:"type"`
	ID                 string     `json:"id"`
	ExternalID         string     `json:"external_id"`
	Released           string     `json:"released"`
	Description        string     `json:"description"`
	URL                string     `json:"url"`
	CollectionURL      string     `json:"collection_url"`
	PreviousVersionURL string     `json:"previous_version_url"`
	RootVersionURL     string     `json:"root_version_url"`
	Extras             Extras     `json:"extras"`
	CreatedOn          time.Time  `json:"created_on"`
	CreatedBy          string     `json:"created_by"`
	UpdatedOn          time.Time  `json:"updated_on"`
	UpdatedBy          string     `json:"updated_by"`
	Collection         Collection `json:"collection"`
}

func (c *Client) CreateCollectionVersion(
	ctx context.Context, collection *Collection, headers *Headers,
) (*Collection, error) {
	var resp Collection

	err := c.makeRequest(ctx, http.MethodPost, composeCollectionVersionsURL(headers), nil, collection, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// composeCollectionVersionsURL forms the create/get collection versions url. It follows this path
// /orgs/:org/collections/:collection/versions/.
func composeCollectionVersionsURL(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/collections/" + headers.Collection + "/versions/"
}
