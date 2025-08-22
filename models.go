package ocl

import (
	"errors"
	"time"
)

type ReleaseVersion struct {
	Released    string `json:"released,omitempty"`
	Description string `json:"description,omitempty"`
}

type CreateVersion struct {
	ID              string `json:"id,omitempty"`
	ExternalID      string `json:"external_id,omitempty"`
	Released        string `json:"released,omitempty"`
	Description     string `json:"description,omitempty"`
	Release         bool   `json:"release,omitempty"`
	PreviousVersion string `json:"previous_version,omitempty"`
	ParentVersion   string `json:"parent_version,omitempty"`
}

type ResourceVersion struct {
	Type               string     `json:"type,omitempty"`
	ID                 string     `json:"id,omitempty"`
	ExternalID         string     `json:"external_id,omitempty"`
	Released           string     `json:"released,omitempty"`
	Description        string     `json:"description,omitempty"`
	URL                string     `json:"url,omitempty"`
	CollectionURL      string     `json:"collection_url,omitempty"`
	PreviousVersionURL string     `json:"previous_version_url,omitempty"`
	RootVersionURL     string     `json:"root_version_url,omitempty"`
	Extras             Extras     `json:"extras"`
	CreatedOn          time.Time  `json:"created_on,omitempty"`
	CreatedBy          string     `json:"created_by,omitempty"`
	UpdatedOn          time.Time  `json:"updated_on,omitempty"`
	UpdatedBy          string     `json:"updated_by,omitempty"`
	Collection         Collection `json:"collection,omitempty"`
}

type ResourceOperationTypeEnum string

const (
	CreateCollectionOperation ResourceOperationTypeEnum = "CREATE_COLLECTION"
	DeleteCollectionOperation ResourceOperationTypeEnum = "DELETE_COLLECTION"
	UpdateCollectionOperation ResourceOperationTypeEnum = "UPDATE_COLLECTION"

	CreateCollectionVersionOperation  ResourceOperationTypeEnum = "CREATE_COLLECTION_VERSION"
	ReleaseCollectionVersionOperation ResourceOperationTypeEnum = "RELEASE_COLLECTION_VERSION"
	RetireCollectionVersionOperation  ResourceOperationTypeEnum = "RETIRE_COLLECTION_VERSION"

	DeleteSourceOrgOperation ResourceOperationTypeEnum = "DELETE_SOURCE"
	UpdateSourceOrgOperation ResourceOperationTypeEnum = "UPDATE_SOURCE"

	CreateSourceVersionOperation  ResourceOperationTypeEnum = "CREATE_SOURCE_VERSION"
	ReleaseSourceVersionOperation ResourceOperationTypeEnum = "RELEASE_SOURCE_VERSION"
	RetireSourceVersionOperation  ResourceOperationTypeEnum = "RETIRE_SOURCE_VERSION"
)

var ErrInvalidIdentifierInput = errors.New(
	"invalid input identifiers: required IDs missing for operation",
)

var ErrReleaseFailure = errors.New(
	"invalid payload: version needs to reviewed before release",
)

// isValidInput checks whether the required identifiers are provided
// for a given resource operation type.
//
// Parameters:
//   - orgID: pointer to the organization ID (may be nil if not required by the operation).
//   - srcID: pointer to the source ID (may be nil if not required).
//   - collectionID: pointer to the collection ID (may be nil if not required).
//   - versionID: pointer to the version ID (may be nil if not required).
//   - operation: the resource operation type being validated (create, update, retire).
//
// Returns:
//   - true if the input parameters are valid for the given operation type.
//   - false otherwise.
//
// Example:
//
//	ok := isValidInput(&orgID, &srcID, &collectionID, &versionID, ReleaseSourceVersionOperation)
//	if !ok {
//	    return errors.New("invalid input for release operation")
//	}
func isValidInput(orgID, srcID, collectionID, versionID *string, operation ResourceOperationTypeEnum) bool {
	switch operation {
	case DeleteSourceOrgOperation, UpdateSourceOrgOperation:
		if orgID == nil || srcID == nil {
			return false
		}
	case UpdateCollectionOperation, DeleteCollectionOperation:
		if orgID == nil || collectionID == nil {
			return false
		}
	case CreateCollectionOperation:
		if orgID == nil {
			return false
		}
	case CreateCollectionVersionOperation, ReleaseCollectionVersionOperation, RetireCollectionVersionOperation:
		if orgID == nil || collectionID == nil {
			return false
		}
	case CreateSourceVersionOperation, ReleaseSourceVersionOperation, RetireSourceVersionOperation:
		if orgID == nil || srcID == nil || versionID == nil {
			return false
		}
	default:
		return false
	}
	return true
}
