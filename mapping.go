package ocl

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Mappings struct {
	Extras                  Extras    `json:"extras,omitempty"`
	Checksums               Checksums `json:"checksums,omitempty"`
	ExternalID              any       `json:"external_id,omitempty"`
	Retired                 bool      `json:"retired,omitempty"`
	MapType                 string    `json:"map_type,omitempty"`
	Source                  string    `json:"source,omitempty"`
	Owner                   string    `json:"owner,omitempty"`
	OwnerType               string    `json:"owner_type,omitempty"`
	FromConceptCode         string    `json:"from_concept_code,omitempty"`
	FromConceptName         any       `json:"from_concept_name,omitempty"`
	FromConceptURL          string    `json:"from_concept_url,omitempty"`
	ToConceptCode           string    `json:"to_concept_code,omitempty"`
	ToConceptName           any       `json:"to_concept_name,omitempty"`
	ToConceptURL            string    `json:"to_concept_url,omitempty"`
	FromSourceOwner         string    `json:"from_source_owner,omitempty"`
	FromSourceOwnerType     string    `json:"from_source_owner_type,omitempty"`
	FromSourceURL           string    `json:"from_source_url,omitempty"`
	FromSourceName          string    `json:"from_source_name,omitempty"`
	ToSourceOwner           string    `json:"to_source_owner,omitempty"`
	ToSourceOwnerType       string    `json:"to_source_owner_type,omitempty"`
	ToSourceURL             string    `json:"to_source_url,omitempty"`
	ToSourceName            string    `json:"to_source_name,omitempty"`
	URL                     string    `json:"url,omitempty"`
	Version                 string    `json:"version,omitempty"`
	ID                      string    `json:"id,omitempty"`
	VersionedObjectID       int       `json:"versioned_object_id,omitempty"`
	VersionedObjectURL      string    `json:"versioned_object_url,omitempty"`
	IsLatestVersion         bool      `json:"is_latest_version,omitempty"`
	UpdateComment           any       `json:"update_comment,omitempty"`
	VersionURL              string    `json:"version_url,omitempty"`
	UUID                    string    `json:"uuid,omitempty"`
	VersionCreatedOn        time.Time `json:"version_created_on,omitempty"`
	FromSourceVersion       any       `json:"from_source_version,omitempty"`
	ToSourceVersion         any       `json:"to_source_version,omitempty"`
	FromConceptNameResolved string    `json:"from_concept_name_resolved,omitempty"`
	ToConceptNameResolved   string    `json:"to_concept_name_resolved,omitempty"`
	Type                    string    `json:"type,omitempty"`
	SortWeight              any       `json:"sort_weight,omitempty"`
	VersionUpdatedOn        time.Time `json:"version_updated_on,omitempty"`
	VersionUpdatedBy        string    `json:"version_updated_by,omitempty"`
	LatestSourceVersion     any       `json:"latest_source_version,omitempty"`
	CreatedOn               time.Time `json:"created_on,omitempty"`
	UpdatedOn               time.Time `json:"updated_on,omitempty"`
	CreatedBy               string    `json:"created_by,omitempty"`
	UpdatedBy               string    `json:"updated_by,omitempty"`
	PublicCanView           bool      `json:"public_can_view,omitempty"`
}

type MappingsInput struct {
	FromConceptURL string `json:"from_concept_url" validate:"required"`
	MapType        string `json:"map_type" validate:"required"`
	ToConceptURL   string `json:"to_concept_url" validate:"required"`
}

type Checksums struct {
	Standard string `json:"standard,omitempty"`
	Smart    string `json:"smart,omitempty"`
}

func (m *MappingsInput) constructConceptURL(organization, source, conceptID string) string {
	return fmt.Sprintf("/orgs/%s/sources/%s/concepts/%s/", organization, source, conceptID)
}

func (m *MappingsInput) ConstructFromConceptURL(organization, source, conceptID string) string {
	return m.constructConceptURL(organization, source, conceptID)
}

func (m *MappingsInput) ConstructToConceptURL(organization, source, conceptID string) string {
	return m.constructConceptURL(organization, source, conceptID)
}

func (c *Client) CreateMappings(ctx context.Context, mappings *MappingsInput, headers *Headers) (*Mappings, error) {
	var resp Mappings

	err := c.makeRequest(ctx, http.MethodPost, composeMappingsPath(headers), nil, mappings, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
