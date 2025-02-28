package ocl

import (
	"context"
	"net/http"
	"time"
)

type CollectionVersion struct {
	Type               string     `json:"type,omitempty"`
	ID                 string     `json:"id,omitempty"`
	ExternalID         string     `json:"external_id,omitempty"`
	Released           string     `json:"released,omitempty"`
	Description        string     `json:"description,omitempty"`
	URL                string     `json:"url,omitempty"`
	CollectionURL      string     `json:"collection_url,omitempty"`
	PreviousVersionURL string     `json:"previous_version_url,omitempty"`
	RootVersionURL     string     `json:"root_version_url,omitempty"`
	Extras             Extras     `json:"extras,omitempty"`
	CreatedOn          time.Time  `json:"created_on,omitempty"`
	CreatedBy          string     `json:"created_by,omitempty"`
	UpdatedOn          time.Time  `json:"updated_on,omitempty"`
	UpdatedBy          string     `json:"updated_by,omitempty"`
	Collection         Collection `json:"collection,omitempty"`
}

func (c *Client) CreateCollectionVersion(
	ctx context.Context, collection *Collection, headers *Headers,
) (*Collection, error) {
	var resp *Collection

	err := c.makeRequest(ctx, http.MethodPost, composeCollectionVersionsURL(headers), nil, collection, &resp)
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
