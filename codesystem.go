package ocl

import (
	"context"
	"fmt"
	"net/http"

	model "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
)

func (c *Client) GetCodeSytem(ctx context.Context, headers Headers) (*model.CodeSystem, error) {
	codeSystem := model.CodeSystem{}

	url := fmt.Sprintf("orgs/%s/CodeSystem/%s", headers.Organisation, headers.Source)

	err := c.makeRequest(ctx, http.MethodGet, url, nil, nil, &codeSystem)
	if err != nil {
		return nil, err
	}

	return &codeSystem, nil
}

func (c *Client) GetAllCodeSystem(ctx context.Context, headers Headers) (*model.Bundle, error) {
	output := model.Bundle{}

	err := c.makeRequest(ctx, http.MethodGet, "fhir/CodeSystem/", nil, nil, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CodeSystem %w", err)
	}

	return &output, nil
}
