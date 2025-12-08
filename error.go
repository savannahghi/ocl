package ocl

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// APIErrorResponse represents the structure that an error message will be returned with.
type APIErrorResponse struct {
	All []string `json:"__all__"`
}

// APIError represents a structured API error response.
type APIError struct {
	StatusCode int
	RawBody    string
	Mnemonic   string `json:"mnemonic"`
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

// IsDuplicateCollectionIDError checks if an error is due to a duplicate Collection  ID within a source.
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

// IsDuplicateSourceIDError checks if an error is due to a duplicate Source ID within an organization.
func IsDuplicateSourceIDError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusBadRequest {
			if slices.Contains(apiErr.APIError.All, "Constraint “org_source_unique” is violated.") {
				return true
			}
		}
	}

	return false
}

// IsDuplicateMappingError checks if an error is due to a duplicate Source ID within an organization.
func IsDuplicateMappingError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusBadRequest {
			if slices.Contains(
				apiErr.APIError.All,
				"Mapping ID must be unique within a source.",
			) {
				return true
			}
		}
	}

	return false
}

func IsDuplicateMnemonicError(err error) bool { //nolint:cyclop
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	hasDupWords := func(s string) bool {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return false
		}

		if strings.Contains(s, "mnemonic") &&
			(strings.Contains(s, "already exists") ||
				strings.Contains(s, "must be unique") ||
				strings.Contains(s, "to_concept_code must be unique.") ||
				strings.Contains(s, "unique")) {
			return true
		}

		return false
	}

	if hasDupWords(apiErr.Mnemonic) {
		return true
	}

	if slices.ContainsFunc(apiErr.APIError.All, hasDupWords) {
		return true
	}

	if hasDupWords(apiErr.RawBody) {
		return true
	}

	return false
}
