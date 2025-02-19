package ocl

import (
	"errors"
	"fmt"
	"net/http"
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
			for _, msg := range apiErr.APIError.All {
				if msg == "Concept ID must be unique within a source." {
					return true
				}
			}
		}
	}

	return false
}
