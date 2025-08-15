package ocl

import (
	"context"
	"fmt"
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

// This input has been informed by the following resouces
type CollectionUpdateInput struct {
	ExternalID       *string   `json:"external_id,omitempty"`
	Name             *string   `json:"name,omitempty"`
	FullName         *string   `json:"full_name,omitempty"`
	CollectionType   *string   `json:"collection_type,omitempty"`
	PublicAccess     *string   `json:"public_access,omitempty"`
	SupportedLocales []*string `json:"supported_locales,omitempty"`
	Website          *string   `json:"website,omitempty"`
	Description      *string   `json:"description,omitempty"`
	Extras           *Extras   `json:"extras,omitempty"`
}

type Extras struct{}

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

func (c *Client) CreateCollection(ctx context.Context, collection *Collection, headers *Headers) (*Collection, error) {
	var resp Collection

	path := fmt.Sprintf("orgs/%s/collections/", headers.Organisation)
	err := c.makeRequest(ctx, http.MethodPost, path, nil, collection, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) RetireCollection(ctx context.Context, headers *Headers, id string) error {
	path := fmt.Sprintf("orgs/%s/collections/%s/", headers.Organisation, id)

	err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateCollection(ctx context.Context, input *CollectionUpdateInput, organizationID, collectionID string) (*Collection, error) {
	var output Collection

	if organizationID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID cannot be empty")
	}

	path := fmt.Sprintf("orgs/%s/collections/%s/", organizationID, collectionID)

	if err := c.makeRequest(ctx, http.MethodPut, path, nil, input, &output); err != nil {
		return nil, err
	}
	return &output, nil
}
