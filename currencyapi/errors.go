package currencyapi

import (
	"errors"
	"fmt"
)

// ValidationError represents an input validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// RequestError represents an error that occurred during request preparation or execution
type RequestError struct {
	Op  string // Operation that failed (e.g., "create_request", "execute_request")
	Err error  // Underlying error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("request error during %s: %v", e.Op, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// HTTPError represents an unexpected HTTP response
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("unexpected HTTP status %d: %s", e.StatusCode, e.Body)
}

// IsNotFound returns true if the error is a 404 Not Found error
func (e *HTTPError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error
func (e *HTTPError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error
func (e *HTTPError) IsRateLimited() bool {
	return e.StatusCode == 429
}

// APIError represents an error returned by the CurrencyAPI
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Info       string
}

func (e *APIError) Error() string {
	if e.Info != "" {
		return fmt.Sprintf("API error [%s]: %s (%s)", e.Code, e.Message, e.Info)
	}
	return fmt.Sprintf("API error [%s]: %s", e.Code, e.Message)
}

// IsQuotaExceeded returns true if the API quota has been exceeded
func (e *APIError) IsQuotaExceeded() bool {
	return e.Code == "quota_exceeded"
}

// IsInvalidAPIKey returns true if the API key is invalid
func (e *APIError) IsInvalidAPIKey() bool {
	return e.Code == "invalid_api_key" || e.StatusCode == 401
}

// ParseError represents an error that occurred while parsing an API response
type ParseError struct {
	Endpoint string
	Err      error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse response from '%s': %v", e.Endpoint, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// APIErrorResponse represents the error response structure from the API
type APIErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Info    string `json:"info,omitempty"`
	} `json:"error"`
}

// Helper functions for error type checking

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// IsRequestError checks if the error is a request error
func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

// IsHTTPError checks if the error is an HTTP error
func IsHTTPError(err error) bool {
	var he *HTTPError
	return errors.As(err, &he)
}

// IsAPIError checks if the error is an API error
func IsAPIError(err error) bool {
	var ae *APIError
	return errors.As(err, &ae)
}

// IsParseError checks if the error is a parse error
func IsParseError(err error) bool {
	var pe *ParseError
	return errors.As(err, &pe)
}

// GetHTTPStatusCode extracts the HTTP status code from an error if available
func GetHTTPStatusCode(err error) (int, bool) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode, true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode, true
	}
	return 0, false
}

// IsTemporaryError returns true if the error might be resolved by retrying
func IsTemporaryError(err error) bool {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		// 5xx errors and 429 are typically temporary
		return httpErr.StatusCode >= 500 || httpErr.StatusCode == 429
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500 || apiErr.StatusCode == 429
	}

	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		// Network errors are typically temporary
		return reqErr.Op == "execute_request"
	}

	return false
}
