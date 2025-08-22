package ocl

import (
	"context"
	"fmt"
	"net/http"
)

// CreateCollectionVersion makes a POST request to
// /orgs/:org/collections/:collection/versions/:version/
// to create a new collection version in OCL
//
// Parameters:
//   - Organization: the ID of the organization (This is a required field or the operation fails).
//   - Collection: the ID of the collection (This is a required field or the operation fails).
//   - CreateVersion: this is the payload containing the fields (This is a required field or the operation fails).
//     to create collection version in OCL.
//
// Returns:
//   - ResourceVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) CreateCollectionVersion(
	ctx context.Context, input *CreateVersion, headers *Headers,
) (*ResourceVersion, error) {
	var resp *ResourceVersion
	if !isValidInput(&headers.Organisation, nil, &headers.Collection, nil, CreateCollectionVersionOperation) {
		return nil, ErrInvalidIdentifierInput
	}

	if input.Release {
		return nil, ErrReleaseFailure
	}
	path := fmt.Sprintf("orgs/%s/collections/%s/versions/", headers.Organisation, headers.Collection)
	err := c.makeRequest(ctx, http.MethodPost, path, nil, input, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ReleaseCollectionVersion makes a POST request to
// /orgs/:org/collections/:collection/versions/:version/
// to make a collection version release on OCL
//
// Parameters:
//   - Organization: the ID of the organization (This is a required field or the operation fails).
//   - Collection: the ID of the collection (This is a required field or the operation fails).
//   - VersionID: the collection version ID (This is a required field or the operation fails).
//   - ReleaseVersion: this is the payload containing the fields to release the collection version in OCL.
//
// Returns:
//   - ResourceVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) ReleaseCollectionVersion(ctx context.Context, headers *Headers, input *ReleaseVersion) (*ResourceVersion, error) {
	if !isValidInput(&headers.Organisation, nil, &headers.Collection, &headers.VersionID, ReleaseCollectionVersionOperation) {
		return nil, ErrInvalidIdentifierInput
	}

	var output *ResourceVersion
	path := fmt.Sprintf("/orgs/%s/collections/%s/versions/%s/", headers.Organisation, headers.Collection, headers.VersionID)
	if err := c.makeRequest(ctx, http.MethodPost, path, nil, input, output); err != nil {
		return nil, err
	}
	return output, nil
}

// RetireCollectionVersion makes a DELETE request to
// /orgs/:org/collections/:collection/versions/:version/
// to retire a collection version in OCL
//
// Parameters:
//   - Organization: the ID of the organization (This is a required field or the operation fails).
//   - Collection: the ID of the collection (This is a required field or the operation fails).
//   - VersionID: the collection version ID (This is a required field or the operation fails).
//
// Returns:
//   - error if the operation fails
//   - no error if operation succeeds
func (c *Client) RetireCollectionVersion(ctx context.Context, headers *Headers) error {
	if !isValidInput(&headers.Organisation, nil, &headers.Collection, &headers.VersionID, ReleaseCollectionVersionOperation) {
		return ErrInvalidIdentifierInput
	}
	path := fmt.Sprintf("/orgs/%s/collections/%s/versions/%s/", headers.Organisation, headers.Collection, headers.VersionID)
	if err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil, nil); err != nil {
		return err
	}
	return nil
}
