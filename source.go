package ocl

import (
	"context"
	"net/http"
	"time"
)

type Source struct {
	Type                   string    `json:"type,omitempty"`
	UUID                   string    `json:"uuid,omitempty"`
	ID                     string    `json:"id,omitempty"`
	ShortCode              string    `json:"short_code,omitempty"`
	Name                   string    `json:"name,omitempty"`
	FullName               string    `json:"full_name,omitempty"`
	Description            string    `json:"description,omitempty"`
	SourceType             string    `json:"source_type,omitempty"`
	CustomValidationSchema string    `json:"custom_validation_schema,omitempty"`
	PublicAccess           string    `json:"public_access,omitempty"`
	DefaultLocale          string    `json:"default_locale,omitempty"`
	SupportedLocales       []string  `json:"supported_locales,omitempty"`
	Website                string    `json:"website,omitempty"`
	URL                    string    `json:"url,omitempty"`
	Owner                  string    `json:"owner,omitempty"`
	OwnerType              string    `json:"owner_type,omitempty"`
	OwnerURL               string    `json:"owner_url,omitempty"`
	CreatedOn              time.Time `json:"created_on,omitempty"`
	UpdatedOn              time.Time `json:"updated_on,omitempty"`
	CreatedBy              string    `json:"created_by,omitempty"`
	UpdatedBy              string    `json:"updated_by,omitempty"`
	Extras                 Extras    `json:"extras,omitempty"`
	ExternalID             any       `json:"external_id,omitempty"`
	VersionsURL            string    `json:"versions_url,omitempty"`
	Version                string    `json:"version,omitempty"`
	ConceptsURL            string    `json:"concepts_url,omitempty"`
	MappingsURL            string    `json:"mappings_url,omitempty"`
	CanonicalURL           any       `json:"canonical_url,omitempty"`
	Publisher              any       `json:"publisher,omitempty"`
	Purpose                any       `json:"purpose,omitempty"`
	Copyright              any       `json:"copyright,omitempty"`
	ContentType            any       `json:"content_type,omitempty"`
	RevisionDate           any       `json:"revision_date,omitempty"`
	LogoURL                any       `json:"logo_url,omitempty"`
	Text                   any       `json:"text,omitempty"`
	ClientConfigs          []any     `json:"client_configs,omitempty"`
	Experimental           any       `json:"experimental,omitempty"`
	CaseSensitive          any       `json:"case_sensitive,omitempty"`
	CollectionReference    any       `json:"collection_reference,omitempty"`
	HierarchyMeaning       any       `json:"hierarchy_meaning,omitempty"`
	Compositional          any       `json:"compositional,omitempty"`
	VersionNeeded          any       `json:"version_needed,omitempty"`
	HierarchyRootURL       any       `json:"hierarchy_root_url,omitempty"`
	Meta                   any       `json:"meta,omitempty"`
}

func (c *Client) CreateSource(ctx context.Context, source *Source) (*Source, error) {
	var resp Source

	err := c.makeRequest(ctx, http.MethodPost, "orgs/", nil, source, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
