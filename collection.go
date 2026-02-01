package ocl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	Extras           Extras    `json:"extras,omitzero"`
	Owner            string    `json:"owner,omitempty"`
	OwnerType        string    `json:"owner_type,omitempty"`
	OwnerURL         string    `json:"owner_url,omitempty"`
	URL              string    `json:"url,omitempty"`
	VersionsURL      string    `json:"versions_url,omitempty"`
	ConceptsURL      string    `json:"concepts_url,omitempty"`
	MappingsURL      string    `json:"mappings_url,omitempty"`
	Versions         int       `json:"versions,omitempty"`
	CreatedOn        time.Time `json:"created_on,omitzero"`
	CreatedBy        string    `json:"created_by,omitempty"`
	UpdatedOn        time.Time `json:"updated_on,omitzero"`
	UpdatedBy        string    `json:"updated_by,omitempty"`
	Released         bool      `json:"released,omitempty"`
	CanonicalURL     string    `json:"canonical_url,omitempty"`
}

type CollectionInput struct {
	ID               string   `json:"id,omitempty"`
	ShortCode        string   `json:"short_code,omitempty"`
	ExternalID       string   `json:"external_id,omitempty"`
	Name             string   `json:"name,omitempty"`
	FullName         string   `json:"full_name,omitempty"`
	CollectionType   string   `json:"collection_type,omitempty"`
	PublicAccess     string   `json:"public_access,omitempty"`
	PreferredSource  string   `json:"preferred_source,omitempty"`
	SupportedLocales []string `json:"supported_locales,omitempty"`
	Website          string   `json:"website,omitempty"`
	Description      string   `json:"description,omitempty"`
	Extras           Extras   `json:"extras,omitzero"`
	CanonicalURL     string   `json:"canonical_url,omitempty"`
}

type Extras struct{}

func (c *Client) CreateCollection(
	ctx context.Context,
	collection *CollectionInput,
	headers *Headers,
) (*Collection, error) {
	var resp Collection

	params := RequestParameters{
		OrganisationID: &headers.Organisation,
	}
	if !isValidInput(params, CreateCollectionOperation) {
		return nil, ErrInvalidIdentifierInput
	}

	path := fmt.Sprintf("orgs/%s/collections/", headers.Organisation)

	err := c.makeRequest(ctx, http.MethodPost, path, nil, collection, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) RetireCollection(ctx context.Context, headers *Headers) error {
	path := fmt.Sprintf("orgs/%s/collections/%s/", headers.Organisation, headers.Collection)

	err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateCollection(ctx context.Context, input *CollectionInput, headers *Headers) (*Collection, error) {
	var output Collection

	params := RequestParameters{
		OrganisationID: &headers.Organisation,
		CollectionID:   &headers.Collection,
	}
	if !isValidInput(params, UpdateCollectionOperation) {
		return nil, ErrInvalidIdentifierInput
	}

	path := fmt.Sprintf("orgs/%s/collections/%s/", headers.Organisation, headers.Collection)

	err := c.makeRequest(ctx, http.MethodPut, path, nil, input, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) GetCollection(ctx context.Context, headers *Headers) (*Collection, error) {
	var output Collection

	path := fmt.Sprintf("orgs/%s/collections/%s/", headers.Organisation, headers.Collection)

	err := c.makeRequest(ctx, http.MethodGet, path, nil, nil, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collection %w", err)
	}

	return &output, nil
}

// CollectionExists checks whether a collection exists in OCL and returns a boolean.
func (c *Client) CollectionExists(ctx context.Context, headers *Headers) (bool, error) {
	if headers.Organisation == "" {
		return false, errors.New("organization ID cannot be null")
	}

	if headers.Collection == "" {
		return false, errors.New("collection ID cannot be null")
	}

	_, err := c.GetCollection(ctx, headers)
	if err != nil {
		if ResourceNotFoundErr(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// ListCollectionConcepts lists all concepts referenced in a collection.
// API: GET /orgs/:org/collections/:collection/[:version/]concepts/.
func (c *Client) ListCollectionConcepts(
	ctx context.Context,
	headers *Headers,
	params url.Values,
) ([]*Concept, error) {
	err := validateCollectionHeaders(headers)
	if err != nil {
		return nil, err
	}

	var resp []*Concept

	path := composeCollectionPath(headers) + "/concepts/"

	err = c.makeRequest(ctx, http.MethodGet, path, params, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list collection concepts: %w", err)
	}

	return resp, nil
}

// Helper to validate collection-related headers.
func validateCollectionHeaders(headers *Headers) error {
	if headers == nil {
		return errors.New("headers cannot be nil")
	}

	if headers.Organisation == "" {
		return errors.New("organisation is required")
	}

	if headers.Collection == "" {
		return errors.New("collection is required")
	}

	return nil
}

// composeCollectionPath builds the base path for collection operations
// Pattern: /orgs/{org}/collections/{collection} or /orgs/{org}/collections/{collection}/{version}.
func composeCollectionPath(headers *Headers) string {
	path := "orgs/" + headers.Organisation + "/collections/" + headers.Collection

	if headers.VersionID != "" {
		path += "/" + headers.VersionID
	}

	return path
}

// CollectionConceptSearchParams contains parameters for searching concepts in a collection expansion.
type CollectionConceptSearchParams struct {
	Query  string
	Limit  int
	Offset int
}

// SearchCollectionConcepts searches for concepts within an OCL collection expansion.
// This calls the endpoint: GET /orgs/:org/collections/:collection/HEAD/expansions/autoexpand-HEAD/concepts/.
func (c *Client) SearchCollectionConcepts(
	ctx context.Context,
	headers *Headers,
	params CollectionConceptSearchParams,
) ([]Concept, error) {
	if params.Limit <= 0 {
		params.Limit = 25
	}

	query := url.Values{}
	if params.Query != "" {
		query.Set("q", params.Query)
	}

	query.Set("limit", strconv.Itoa(params.Limit))
	query.Set("offset", strconv.Itoa(params.Offset))

	var resp []Concept

	err := c.makeRequest(
		ctx,
		http.MethodGet,
		composeCollectionExpansionConceptsURL(headers),
		query,
		nil,
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search collection concepts: %w", err)
	}

	return resp, nil
}

// composeCollectionExpansionConceptsURL forms the URL to search concepts in a collection expansion.
// It follows this path: /orgs/{org}/collections/{collection}/HEAD/expansions/autoexpand-HEAD/concepts/.
func composeCollectionExpansionConceptsURL(headers *Headers) string {
	return fmt.Sprintf(
		"orgs/%s/collections/%s/HEAD/expansions/autoexpand-HEAD/concepts/",
		headers.Organisation,
		headers.Collection,
	)
}
