package ocl

import (
	"errors"
	"reflect"
	"time"
)

type ReleaseVersion struct {
	Released    string `json:"released,omitempty"`
	Description string `json:"description,omitempty"`
}

type CollectionVersionInput struct {
	ID           string `json:"id,omitempty"`
	Released     bool   `json:"released,omitempty"`
	Description  string `json:"description,omitempty"`
	Release      bool   `json:"release,omitempty"`
	ExpansionURL string `json:"expansion_url,omitempty"`
	AutoExapand  bool   `json:"autoexpand,omitempty"`
}

type CollectionVersion struct {
	Type               string     `json:"type,omitempty"`
	ID                 string     `json:"id,omitempty"`
	ExternalID         string     `json:"external_id,omitempty"`
	Released           bool       `json:"released,omitempty"`
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
	Collection         Collection `json:"collection"`
}

type SourceVersionInput struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Released    bool   `json:"released,omitempty"`
}

type SourceVersion struct {
	Type               string    `json:"type,omitempty"`
	ID                 string    `json:"id,omitempty"`
	ExternalID         string    `json:"external_id,omitempty"`
	Released           string    `json:"released,omitempty"`
	Description        string    `json:"description,omitempty"`
	URL                string    `json:"url,omitempty"`
	CollectionURL      string    `json:"collection_url,omitempty"`
	PreviousVersionURL string    `json:"previous_version_url,omitempty"`
	RootVersionURL     string    `json:"root_version_url,omitempty"`
	Extras             Extras    `json:"extras"`
	CreatedOn          time.Time `json:"created_on,omitempty"`
	CreatedBy          string    `json:"created_by,omitempty"`
	UpdatedOn          time.Time `json:"updated_on,omitempty"`
	UpdatedBy          string    `json:"updated_by,omitempty"`
	Collection         Source    `json:"collection"`
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

// RequestParameters is a single struct to hold all possible input
type RequestParameters struct {
	OrganisationID *string
	SourceID       *string
	CollectionID   *string
	VersionID      *string
}

// A map to define which parameters are required for each operation
var requiredParams = map[ResourceOperationTypeEnum][]string{
	DeleteSourceOrgOperation:          {"OrganisationID", "SourceID"},
	UpdateSourceOrgOperation:          {"OrganisationID", "SourceID"},
	UpdateCollectionOperation:         {"OrganisationID", "CollectionID"},
	DeleteCollectionOperation:         {"OrganisationID", "CollectionID"},
	CreateCollectionOperation:         {"OrganisationID"},
	CreateCollectionVersionOperation:  {"OrganisationID", "CollectionID"},
	ReleaseCollectionVersionOperation: {"OrganisationID", "CollectionID"},
	RetireCollectionVersionOperation:  {"OrganisationID", "CollectionID"},
	CreateSourceVersionOperation:      {"OrganisationID", "SourceID", "VersionID"},
	ReleaseSourceVersionOperation:     {"OrganisationID", "SourceID", "VersionID"},
	RetireSourceVersionOperation:      {"OrganisationID", "SourceID", "VersionID"},
}

// isValidInput: checks whether the required identifiers for a given operation are not nil
func isValidInput(params RequestParameters, operation ResourceOperationTypeEnum) bool {
	required, ok := requiredParams[operation]
	if !ok {
		return false
	}
	v := reflect.ValueOf(params)
	for _, fieldName := range required {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}
		if field.Kind() == reflect.Pointer && field.IsNil() {
			return false
		}
	}
	return true
}
