package ocl

import (
	"context"
	"fmt"
	model "github.com/savannahghi/hapi-fhir-go/models/r5/fhir500"
	"net/http"
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
