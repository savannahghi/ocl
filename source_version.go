package ocl

import (
	"context"
	"fmt"
	"net/http"
)

// CreateSourceVersion makes a POST request to
// /orgs/:org/sources/:source/versions/
// to create a new source version in OCL
//
// Parameters:
//   - Headers: this contains organization and source IDs.
//   - SourceVersionInput: this is the payload containing the fields to create source version in OCL.
//
// Returns:
//   - SourceVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) CreateSourceVersion(ctx context.Context, headers *Headers, input *SourceVersionInput) (*SourceVersion, error) {
	params := RequestParameters{
		OrganisationID: &headers.Organisation,
		SourceID:       &headers.Source,
	}
	if !isValidInput(params, CreateSourceVersionOperation) {
		return nil, ErrInvalidIdentifierInput
	}

	var output *SourceVersion
	path := fmt.Sprintf("orgs/%s/sources/%s/versions/", headers.Organisation, headers.Source)

	err := c.makeRequest(ctx, http.MethodPost, path, nil, input, &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// ReleaseSourceVersion makes a POST request to
// /orgs/:org/sources/:source/:version/
// to release source version on OCL
//
// Parameters:
//   - Headers: Contains organizatio, source and version IDs.
//   - ReleaseVersion: This is the payload containing the fields to release source vesion on OCL.
//
// Returns:
//   - SourceVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) ReleaseSourceVersion(ctx context.Context, headers *Headers, input ReleaseVersion) (*SourceVersion, error) {
	params := RequestParameters{
		OrganisationID: &headers.Organisation,
		SourceID:       &headers.Source,
		VersionID:      &headers.VersionID,
	}
	if !isValidInput(params, ReleaseSourceVersionOperation) {
		return nil, ErrInvalidIdentifierInput
	}
	var output *SourceVersion

	path := fmt.Sprintf("orgs/%s/sources/%s/%s/", headers.Organisation, headers.Source, headers.VersionID)
	err := c.makeRequest(ctx, http.MethodPost, path, nil, input, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// RetireSourceVersion makes a DELETE request to
// /orgs/:org/sources/:source/:version/
// to deaactivate a source version on OCL
//
// Parameters:
//   - Headers: Contains organization, source and version IDs.
//
// Returns:
//   - nil if the operation succeeds.
//   - error if the operation fails.
func (c *Client) RetireSourceVersion(ctx context.Context, headers *Headers) error {
	params := RequestParameters{
		OrganisationID: &headers.Organisation,
		SourceID:       &headers.Source,
		VersionID:      &headers.VersionID,
	}
	if !isValidInput(params, RetireSourceVersionOperation) {
		return ErrInvalidIdentifierInput
	}
	path := fmt.Sprintf("/orgs/%s/sources/%s/%s/", headers.Organisation, headers.Source, headers.VersionID)
	err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
