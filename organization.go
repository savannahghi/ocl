package ocl

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
)

// SimpleOrganizationInput is a simple model used to create an Organization in advantage.
type SimpleOrganizationInput struct {
	ID           string `json:"id" validate:"required"`
	PublicAccess string `json:"public_access" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Company      string `json:"company" validate:"required"`
	Website      string `json:"website" validate:"required"`
	Location     string `json:"location,omitempty"`
}

var validate = validator.New()

func ValidateStruct(input interface{}) error {
	return validate.Struct(input)
}

type OrganizationOutput struct {
	Type              string    `json:"type,omitempty"`
	UUID              string    `json:"uuid,omitempty"`
	ID                string    `json:"id,omitempty"`
	PublicAccess      string    `json:"public_access,omitempty"`
	Name              string    `json:"name,omitempty"`
	Company           string    `json:"company,omitempty"`
	Website           string    `json:"website,omitempty"`
	Location          string    `json:"location,omitempty"`
	Members           int       `json:"members,omitempty"`
	CreatedOn         time.Time `json:"created_on,omitempty"`
	UpdatedOn         time.Time `json:"updated_on,omitempty"`
	URL               string    `json:"url,omitempty"`
	Extras            any       `json:"extras,omitempty"`
	CreatedBy         string    `json:"created_by,omitempty"`
	UpdatedBy         string    `json:"updated_by,omitempty"`
	SourcesURL        string    `json:"sources_url,omitempty"`
	PublicSources     int       `json:"public_sources,omitempty"`
	CollectionsURL    string    `json:"collections_url,omitempty"`
	PublicCollections int       `json:"public_collections,omitempty"`
	LogoURL           any       `json:"logo_url,omitempty"`
	Description       any       `json:"description,omitempty"`
	Text              any       `json:"text,omitempty"`
}

// CreateOrganization is used to create an organization
func (c *Client) CreateOrganization(ctx context.Context, organization SimpleOrganizationInput, headers *Headers) (*OrganizationOutput, error) {
	if err := ValidateStruct(organization); err != nil {
		return nil, fmt.Errorf("invalid organization payload: %w", err)
	}

	var output OrganizationOutput

	err := c.makeRequest(ctx, http.MethodPost, "orgs/", nil, organization, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
