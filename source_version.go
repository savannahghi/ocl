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
//   - Organization: the ID of the organization (This is a required field or the operation fails).
//   - Source: the ID of the source (This is a required field or the operation fails).
//   - CreateVersion: this is the payload containing the fields
//     to create source version in OCL.
//
// Returns:
//   - ResourceVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) CreateSourceVersion(ctx context.Context, headers *Headers, input *CreateVersion) (*ResourceVersion, error) {
	if !isValidInput(&headers.Organisation, &headers.Source, nil, nil, CreateSourceVersionOperation) {
		return nil, ErrInvalidIdentifierInput
	}

	if input.Release {
		return nil, ErrReleaseFailure
	}

	var output *ResourceVersion
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
//   - Organization: the ID of the organization (This is a required field or the operation fails).
//   - Source: the ID of the source (This is a required field or the operation fails).
//   - VersionID: the version ID of the source (This is a required field or the operation fails)
//   - ReleaseVersion: this is the payload containing the fields to release source vesion on OCL.
//
// Returns:
//   - ResourceVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) ReleaseSourceVersion(ctx context.Context, sourceID, organizationID, versionID string, input ReleaseVersion) (*ResourceVersion, error) {
	if !isValidInput(&organizationID, &sourceID, nil, &versionID, ReleaseSourceVersionOperation) {
		return nil, ErrInvalidIdentifierInput
	}
	var output *ResourceVersion

	path := fmt.Sprintf("orgs/%s/sources/%s/%s/", organizationID, sourceID, versionID)
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
//   - Organization: the ID of the organization (This is a required field or the operation fails).
//   - Source: the ID of the source (This is a required field or the operation fails).
//   - VersionID: the version ID of the source (This is a required field or the operation fails)
//
// Returns:
//   - nil if the operation succeeds.
//   - error if the operation fails.
func (c *Client) RetireSourceVersion(ctx context.Context, organizationID, sourceID, versionID string) error {
	if !isValidInput(&organizationID, &sourceID, nil, &versionID, RetireSourceVersionOperation) {
		return ErrInvalidIdentifierInput
	}
	path := fmt.Sprintf("/orgs/%s/sources/%s/%s/", organizationID, sourceID, versionID)
	err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
