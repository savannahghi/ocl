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
//   - Headers: Contains organization and collection IDs.
//   - CollectionVersionInput: this is the payload containing the fields (This is a required field or the operation fails).
//     to create collection version in OCL.
//
// Returns:
//   - CollectionVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) CreateCollectionVersion(
	ctx context.Context, input *CollectionVersionInput, headers *Headers,
) (*CollectionVersion, error) {
	var resp *CollectionVersion
	params := RequestParameters{
		OrganisationID: &headers.Organisation,
		CollectionID:   &headers.Collection,
	}
	if !isValidInput(params, CreateCollectionVersionOperation) {
		return nil, ErrInvalidIdentifierInput
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
//   - Headers: Contains organization, collection and version IDs.
//   - ReleaseVersion: this is the payload containing the fields to release the collection version in OCL.
//
// Returns:
//   - CollectionVersion if the operation succeeds.
//   - error if the operation fails.
func (c *Client) ReleaseCollectionVersion(ctx context.Context, headers *Headers, input *ReleaseVersion) (*CollectionVersion, error) {
	params := RequestParameters{
		OrganisationID: &headers.Organisation,
		CollectionID:   &headers.Collection,
		VersionID:      &headers.VersionID,
	}
	if !isValidInput(params, ReleaseCollectionVersionOperation) {
		return nil, ErrInvalidIdentifierInput
	}

	var output *CollectionVersion
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
//   - Headers: Contains organization, collection, and version IDs.
//
// Returns:
//   - error if the operation fails
//   - no error if operation succeeds
func (c *Client) RetireCollectionVersion(ctx context.Context, headers *Headers) error {
	params := RequestParameters{
		OrganisationID: &headers.Organisation,
		CollectionID:   &headers.Collection,
		VersionID:      &headers.VersionID,
	}
	if !isValidInput(params, ReleaseCollectionVersionOperation) {
		return ErrInvalidIdentifierInput
	}
	path := fmt.Sprintf("/orgs/%s/collections/%s/versions/%s/", headers.Organisation, headers.Collection, headers.VersionID)
	if err := c.makeRequest(ctx, http.MethodDelete, path, nil, nil, nil); err != nil {
		return err
	}
	return nil
}
