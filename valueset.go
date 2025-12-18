package ocl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	model "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
)

// GetValueSetVersion retrieves a specific version of a ValueSet from OCL.
// Uses the canonical URL and version ID to fetch an exact ValueSet version,
// OCL API endpoint pattern: /fhir/ValueSet/{query_params}/*.
func (c *Client) GetValueSetVersion(
	ctx context.Context, canonicalURL string,
	headers *Headers,
) (*model.Bundle, error) {
	output := model.Bundle{}

	path := "fhir/ValueSet/"
	params := url.Values{}

	params.Add("url", canonicalURL)
	params.Add("version", headers.VersionID)

	err := c.makeRequest(ctx, http.MethodGet, path, params, nil, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ValueSet version %w", err)
	}

	return &output, nil
}
