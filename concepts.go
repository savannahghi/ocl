package ocl

import (
	"context"
	"net/http"
	"time"
)

// Concept is a unit of meaning that can represent a clinical idea e.g a disease, symptom, medication etc.
// Each concept is uniquely identified within the system and has different attributes attached to it.
type Concept struct {
	UUID                string                 `json:"uuid,omitempty"`
	Extras              map[string]interface{} `json:"extras,omitempty"`
	Checksums           string                 `json:"checksums,omitempty"`
	ID                  string                 `json:"id,omitempty"`
	ExternalID          string                 `json:"external_id,omitempty"`
	ConceptClass        string                 `json:"concept_class,omitempty"`
	Datatype            string                 `json:"datatype,omitempty"`
	URL                 string                 `json:"url,omitempty"`
	Retired             bool                   `json:"retired,omitempty"`
	Source              string                 `json:"source,omitempty"`
	Owner               string                 `json:"owner,omitempty"`
	OwnerType           string                 `json:"owner_type,omitempty"`
	OwnerURL            string                 `json:"owner_url,omitempty"`
	DisplayName         string                 `json:"display_name,omitempty"`
	DisplayLocale       string                 `json:"display_locale,omitempty"`
	Names               []Names                `json:"names,omitempty"`
	Descriptions        []Descriptions         `json:"descriptions,omitempty"`
	CreatedOn           time.Time              `json:"created_on,omitempty"`
	UpdatedOn           time.Time              `json:"updated_on,omitempty"`
	VersionsURL         string                 `json:"versions_url,omitempty"`
	Version             string                 `json:"version,omitempty"`
	ParentID            string                 `json:"parent_id,omitempty"`
	Type                string                 `json:"type,omitempty"`
	UpdateComment       string                 `json:"update_comment,omitempty"`
	VersionURL          string                 `json:"version_url,omitempty"`
	UpdatedBy           time.Time              `json:"updated_by,omitempty"`
	CreatedBy           time.Time              `json:"created_by,omitempty"`
	PublicCanView       string                 `json:"public_can_view,omitempty"`
	VersionedObjectID   string                 `json:"versioned_object_id,omitempty"`
	LatestSourceVersion string                 `json:"latest_source_version,omitempty"`
}

type Names struct {
	UUID            string `json:"uuid,omitempty"`
	Name            string `json:"name,omitempty"`
	ExternalID      string `json:"external_id,omitempty"`
	Type            string `json:"type,omitempty"`
	Locale          string `json:"locale,omitempty"`
	LocalePreferred bool   `json:"locale_preferred,omitempty"`
	NameType        string `json:"name_type,omitempty"`
	Checksum        string `json:"checksum,omitempty"`
}

type Descriptions struct {
	UUID            string `json:"uuid,omitempty"`
	Description     string `json:"description,omitempty"`
	ExternalID      string `json:"external_id,omitempty"`
	Type            string `json:"type,omitempty"`
	Locale          string `json:"locale,omitempty"`
	LocalePreferred bool   `json:"locale_preferred,omitempty"`
	DescriptionType string `json:"description_type,omitempty"`
	Checksum        string `json:"checksum,omitempty"`
}

func (c *Client) CreateConcept(ctx context.Context, concept *Concept) (*Concept, error) {
	var resp Concept

	err := c.makeRequest(ctx, http.MethodPost, "concepts/", nil, concept, resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
