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

type Mapping struct {
	Type                string    `json:"type"`
	UUID                string    `json:"uuid"`
	ExternalID          string    `json:"external_id"`
	Retired             bool      `json:"retired"`
	MapType             string    `json:"map_type"`
	FromSourceOwner     string    `json:"from_source_owner"`
	FromSourceOwnerType string    `json:"from_source_owner_type"`
	FromSourceName      string    `json:"from_source_name"`
	FromConceptCode     string    `json:"from_concept_code"`
	FromConceptName     string    `json:"from_concept_name"`
	FromSourceURL       string    `json:"from_source_url"`
	FromConceptURL      string    `json:"from_concept_url"`
	ToSourceOwner       string    `json:"to_source_owner"`
	ToSourceOwnerType   string    `json:"to_source_owner_type"`
	ToSourceName        string    `json:"to_source_name"`
	ToConceptCode       string    `json:"to_concept_code"`
	ToConceptName       string    `json:"to_concept_name"`
	ToSourceURL         string    `json:"to_source_url"`
	Source              string    `json:"source"`
	Owner               string    `json:"owner"`
	OwnerType           string    `json:"owner_type"`
	OwnerURL            string    `json:"owner_url"`
	URL                 string    `json:"url"`
	Extras              struct{}  `json:"extras"`
	CreatedOn           time.Time `json:"created_on"`
	CreatedBy           string    `json:"created_by"`
	UpdatedOn           time.Time `json:"updated_on"`
	UpdatedBy           string    `json:"updated_by"`
}

type MappingInput struct {
	ID              string `json:"id,omitempty"`
	MapType         string `json:"map_type,omitempty"`
	FromSourceURL   string `json:"from_source_url,omitempty"`
	FromConceptCode string `json:"from_concept_code,omitempty"`
	FromConceptName string `json:"from_concept_name,omitempty"`
	ToSourceURL     string `json:"to_source_url,omitempty"`
	ToConceptCode   string `json:"to_concept_code,omitempty"`
	ToConceptName   string `json:"to_concept_name,omitempty"`
	Owner           string `json:"owner,omitempty"`
	Source          string `json:"source,omitempty"`
}

type Checksums struct {
	Standard string `json:"standard,omitempty"`
	Smart    string `json:"smart,omitempty"`
}

func (m *MappingInput) constructConceptURL(organization, source, conceptID string) string {
	return fmt.Sprintf("/orgs/%s/sources/%s/concepts/%s/", organization, source, conceptID)
}

func (m *MappingInput) ConstructFromConceptURL(organization, source, conceptID string) string {
	return m.constructConceptURL(organization, source, conceptID)
}

func (m *MappingInput) ConstructToConceptURL(organization, source, conceptID string) string {
	return m.constructConceptURL(organization, source, conceptID)
}

func (c *Client) CreateMappings(ctx context.Context, mappings *MappingInput, headers *Headers) (*Mapping, error) {
	var resp Mapping

	err := c.makeRequest(ctx, http.MethodPost, composeMappingsPath(headers), nil, mappings, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) UpdateMappings(ctx context.Context, mappings *MappingInput, headers *Headers) (*Mapping, error) {
	var resp Mapping

	err := c.makeRequest(ctx, http.MethodPut, composeUpdateMappingsPath(headers), nil, mappings, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
