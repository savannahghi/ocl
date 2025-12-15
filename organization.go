package ocl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
)

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
	CreatedOn         time.Time `json:"created_on,omitzero"`
	UpdatedOn         time.Time `json:"updated_on,omitzero"`
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

// SimpleOrganizationInput is a simple model used to create an Organization in advantage.
type SimpleOrganizationInput struct {
	ID           string `json:"id"                    validate:"required"`
	PublicAccess string `json:"public_access"         validate:"required"`
	Name         string `json:"name"                  validate:"required"`
	Company      string `json:"company"               validate:"required"`
	Website      string `json:"website"               validate:"required"`
	Location     string `json:"location,omitempty"`
	Extras       any    `json:"extras,omitempty"`
	Description  string `json:"description,omitempty"`
	Text         string `json:"text,omitempty"`
}

var validate = validator.New()

func ValidateStruct(input any) error {
	return validate.Struct(input)
}

// CreateOrganization is used to create an organization.
func (c *Client) CreateOrganization(
	ctx context.Context,
	organization SimpleOrganizationInput,
) (*OrganizationOutput, error) {
	err := ValidateStruct(organization)
	if err != nil {
		return nil, fmt.Errorf("invalid organization payload: %w", err)
	}

	var output OrganizationOutput

	err = c.makeRequest(ctx, http.MethodPost, "orgs/", nil, organization, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) UpdateOrganization(
	ctx context.Context,
	organization SimpleOrganizationInput,
) (*OrganizationOutput, error) {
	if organization.ID == "" {
		return nil, errors.New("organization ID cannot be null")
	}

	var output OrganizationOutput

	updatePath := fmt.Sprintf("orgs/%s/", organization.ID)

	err := c.makeRequest(ctx, http.MethodPut, updatePath, nil, organization, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) GetOrganization(ctx context.Context, organizationID string) (*OrganizationOutput, error) {
	if organizationID == "" {
		return nil, errors.New("organization ID cannot be null")
	}

	var output OrganizationOutput

	orgPath := fmt.Sprintf("orgs/%s/", organizationID)

	err := c.makeRequest(ctx, http.MethodGet, orgPath, nil, nil, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) DeleteOrganization(ctx context.Context, organizationID string) error {
	if organizationID == "" {
		return errors.New("organization ID cannot be null")
	}

	orgPath := fmt.Sprintf("orgs/%s/", organizationID)

	err := c.makeRequest(ctx, http.MethodDelete, orgPath, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
