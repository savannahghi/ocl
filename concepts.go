package ocl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Concept struct {
	UUID                string         `json:"uuid,omitempty"`
	Extras              map[string]any `json:"extras,omitempty"`
	Checksums           map[string]any `json:"checksums,omitempty"`
	ID                  string         `json:"id,omitempty"`
	ExternalID          string         `json:"external_id,omitempty"`
	ConceptClass        string         `json:"concept_class,omitempty"`
	Datatype            string         `json:"datatype,omitempty"`
	URL                 string         `json:"url,omitempty"`
	Retired             bool           `json:"retired,omitempty"`
	Source              string         `json:"source,omitempty"`
	Owner               string         `json:"owner,omitempty"`
	OwnerType           string         `json:"owner_type,omitempty"`
	OwnerURL            string         `json:"owner_url,omitempty"`
	DisplayName         string         `json:"display_name,omitempty"`
	DisplayLocale       string         `json:"display_locale,omitempty"`
	Names               []Names        `json:"names,omitempty"`
	Descriptions        []Descriptions `json:"descriptions,omitempty"`
	CreatedOn           time.Time      `json:"created_on,omitzero"`
	UpdatedOn           time.Time      `json:"updated_on,omitzero"`
	VersionsURL         string         `json:"versions_url,omitempty"`
	Version             string         `json:"version,omitempty"`
	ParentID            string         `json:"parent_id,omitempty"`
	Type                string         `json:"type,omitempty"`
	UpdateComment       string         `json:"update_comment,omitempty"`
	VersionURL          string         `json:"version_url,omitempty"`
	UpdatedBy           string         `json:"updated_by,omitempty"`
	CreatedBy           string         `json:"created_by,omitempty"`
	PublicCanView       bool           `json:"public_can_view,omitempty"`
	VersionedObjectID   int            `json:"versioned_object_id,omitempty"`
	LatestSourceVersion string         `json:"latest_source_version,omitempty"`
}

type Names struct {
	UUID       string `json:"uuid,omitempty"`
	Name       string `json:"name,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Type       string `json:"type,omitempty"`
	Locale     string `json:"locale,omitempty"`
	NameType   string `json:"name_type,omitempty"`
	Checksum   string `json:"checksum,omitempty"`
}

type Descriptions struct {
	UUID            string `json:"uuid,omitempty"`
	Description     string `json:"description,omitempty"`
	ExternalID      string `json:"external_id,omitempty"`
	Type            string `json:"type,omitempty"`
	Locale          string `json:"locale,omitempty"`
	DescriptionType string `json:"description_type,omitempty"`
	Checksum        string `json:"checksum,omitempty"`
}

type ConceptPage struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []*Concept `json:"results"`
}

// Concept is a unit of meaning that can represent a clinical idea e.g a disease, symptom, medication etc.
// Each concept is uniquely identified within the system and has different attributes attached to it.

func (c *Client) CreateConcept(ctx context.Context, concept *Concept, headers *Headers) (*Concept, error) {
	var resp Concept

	err := c.makeRequest(ctx, http.MethodPost, composeConceptsURL(headers), nil, concept, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) UpdateConcept(ctx context.Context, concept *Concept, headers *Headers) (*Concept, error) {
	var resp Concept

	err := c.makeRequest(ctx, http.MethodPatch, composeConceptURL(headers), nil, concept, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) FetchConcept(ctx context.Context, headers *Headers) (*Concept, error) {
	var resp Concept

	err := c.makeRequest(ctx, http.MethodGet, composeConceptURL(headers), nil, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch concept: %w", err)
	}

	return &resp, nil
}

func (c *Client) ListConcepts(ctx context.Context, headers *Headers, params url.Values) (*ConceptPage, error) {
	var resp ConceptPage

	err := c.makeRequest(ctx, http.MethodGet, composeConceptsURL(headers), params, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list concepts: %w", err)
	}

	return &resp, nil
}

// composeConceptsURL forms the create/get concepts url. It follows this path
// /orgs/{org}/sources/{source}/concepts/.
func composeConceptsURL(headers *Headers) string {
	return composeOrgSourcePath(headers) + "/concepts/"
}

// composeConceptURL forms the get detail/update concepts url. It follows this path
// /orgs/{org}/sources/{source}/concepts/{concept_id}/.
func composeConceptURL(headers *Headers) string {
	return composeOrgSourcePath(headers) + "/concepts/" + headers.ConceptID + "/"
}

// composeOrgSourcePath creates a url path with the org & source set. This is because most of the
// APIs in OCL are namespaced to the source and organisation. It will follow this structure
// /orgs/{org}/sources/{source}.
func composeOrgSourcePath(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/sources/" + headers.Source
}

// composeMappingsPath creates a url path with the org & source set.
// It will follow this structure
// /orgs/{org}/sources/{source}/mappings.
func composeMappingsPath(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/sources/" + headers.Source + "/mappings/"
}

// composeUpdateMappingsPath creates a url path with the org, source & mapping set.
// It will follow this structure
// /orgs/{org}/sources/{source}/mappings/{mapping}.
func composeUpdateMappingsPath(headers *Headers) string {
	return "orgs/" + headers.Organisation + "/sources/" + headers.Source + "/mappings/" + headers.MappingID
}
