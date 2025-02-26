package ocl

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
)

// APIErrorResponse represents the structure that an error message will be returned with.
type APIErrorResponse struct {
	All []string `json:"__all__"`
}

// APIError represents a structured API error response.
type APIError struct {
	StatusCode int
	RawBody    string
	APIError   APIErrorResponse
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API request failed (HTTP %d): %s", e.StatusCode, e.RawBody)
}

// IsDuplicateConceptIDError checks if an error is due to a duplicate Concept ID within a source.
func IsDuplicateConceptIDError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusBadRequest {
			if slices.Contains(apiErr.APIError.All, "Concept ID must be unique within a source.") {
				return true
			}
		}
	}

	return false
}

// IsDuplicateCollectionIDError checks if an error is due to a duplicate Collection ID within a source.
func IsDuplicateCollectionIDError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusBadRequest {
			if slices.Contains(apiErr.APIError.All, "Constraint “org_collection_unique” is violated.") {
				return true
			}
		}
	}

	return false
}
