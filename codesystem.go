package ocl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	model "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
)

// GetCodeSystem retrieves a single CodeSystem resource from the OCL FHIR service.
//
// This method fetches a complete CodeSystem including all concepts, properties, and metadata
// from Open Concept Lab (OCL)
//
// OCL API endpoint pattern: /orgs/{org}/CodeSystem/{source}.
func (c *Client) GetCodeSystem(ctx context.Context, headers Headers) (*model.CodeSystem, error) {
	codeSystem := model.CodeSystem{}

	url := fmt.Sprintf("orgs/%s/CodeSystem/%s", headers.Organisation, headers.Source)

	err := c.makeRequest(ctx, http.MethodGet, url, nil, nil, &codeSystem)
	if err != nil {
		return nil, err
	}

	return &codeSystem, nil
}

// GetAllCodeSystems retrieves all CodeSystems from an OCL organization.
//
// Returns a FHIR Bundle containing the latest version of each CodeSystem.
// OCL API endpoint pattern: /orgs/{org}/CodeSystem.
func (c *Client) GetAllCodeSystems(ctx context.Context, headers Headers) (*model.Bundle, error) {
	output := model.Bundle{}

	url := fmt.Sprintf("orgs/%s/CodeSystem", headers.Organisation)

	err := c.makeRequest(ctx, http.MethodGet, url, nil, nil, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CodeSystem %w", err)
	}

	return &output, nil
}

// GetCodeSystemVersion retrieves a specific version of a CodeSystem from OCL.
//
// Uses the canonical URL and version ID to fetch an exact CodeSystem version,
// which is useful for ensuring consistent terminology across deployments.
// OCL API endpoint pattern: /fhir/CodeSystem/{query_params}/*.
func (c *Client) GetCodeSystemVersion(
	ctx context.Context, cannonicalURL string,
	headers *Headers,
) (*model.Bundle, error) {
	output := model.Bundle{}

	path := "fhir/CodeSystem/"
	params := url.Values{}

	params.Add("url", cannonicalURL)
	params.Add("version", headers.VersionID)

	err := c.makeRequest(ctx, http.MethodGet, path, params, nil, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CodeSystem version %w", err)
	}

	return &output, nil
}
