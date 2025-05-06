package ocl

import (
	"context"
	"net/http"
	"time"
)

type Collection struct {
	Type             string    `json:"type,omitempty"`
	UUID             string    `json:"uuid,omitempty"`
	ID               string    `json:"id,omitempty"`
	ExternalID       string    `json:"external_id,omitempty"`
	ShortCode        string    `json:"short_code,omitempty"`
	Name             string    `json:"name,omitempty"`
	FullName         string    `json:"full_name,omitempty"`
	CollectionType   string    `json:"collection_type,omitempty"`
	PublicAccess     string    `json:"public_access,omitempty"`
	SupportedLocales []string  `json:"supported_locales,omitempty"`
	Website          string    `json:"website,omitempty"`
	Description      string    `json:"description,omitempty"`
	PreferredSource  string    `json:"preferred_source,omitempty"`
	Extras           Extras    `json:"extras,omitempty"`
	Owner            string    `json:"owner,omitempty"`
	OwnerType        string    `json:"owner_type,omitempty"`
	OwnerURL         string    `json:"owner_url,omitempty"`
	URL              string    `json:"url,omitempty"`
	VersionsURL      string    `json:"versions_url,omitempty"`
	ConceptsURL      string    `json:"concepts_url,omitempty"`
	MappingsURL      string    `json:"mappings_url,omitempty"`
	Versions         int       `json:"versions,omitempty"`
	CreatedOn        time.Time `json:"created_on,omitempty"`
	CreatedBy        string    `json:"created_by,omitempty"`
	UpdatedOn        time.Time `json:"updated_on,omitempty"`
	UpdatedBy        string    `json:"updated_by,omitempty"`
	Released         bool      `json:"released,omitempty"`
	CanonicalURL     string    `json:"canonical_url,omitempty"`
}

type Extras struct{}

func (c *Client) CreateCollection(ctx context.Context, collection *Collection, headers *Headers) (*Collection, error) {
	var resp Collection

	err := c.makeRequest(ctx, http.MethodPost, composeCollectionsURL(headers), nil, collection, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// composeCollectionsURL forms the create/get collections url. It follows this path
// /orgs/{org}/collections/.
func composeCollectionsURL(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/collections/"
}
