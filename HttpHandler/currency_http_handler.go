package HttpHandler

import (
	"Golang/currency_converter"
	"Golang/currencyapi"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func CurrencyStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status, err := currency_converter.CheckStatus()
	if err != nil {
		handleCurrencyError(w, err)
		return
	}

	fmt.Fprintf(w, "%s", status)
}

// handleCurrencyError handles errors from the currency converter with appropriate HTTP responses
func handleCurrencyError(w http.ResponseWriter, err error) {
	log.Printf("Currency API error: %v", err)

	// Check for specific error types and set appropriate status codes
	var apiErr *currencyapi.APIError
	if errors.As(err, &apiErr) {
		if apiErr.IsInvalidAPIKey() {
			http.Error(w, `{"error": "Service configuration error"}`, http.StatusInternalServerError)
			return
		}
		if apiErr.IsQuotaExceeded() {
			http.Error(w, `{"error": "Service temporarily unavailable, please try again later"}`, http.StatusServiceUnavailable)
			return
		}
	}

	var httpErr *currencyapi.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.IsRateLimited() {
			http.Error(w, `{"error": "Rate limited, please try again later"}`, http.StatusTooManyRequests)
			return
		}
	}

	// Check if it's a temporary error
	if currencyapi.IsTemporaryError(err) {
		http.Error(w, `{"error": "Service temporarily unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	// Default error response
	http.Error(w, `{"error": "Failed to fetch currency data"}`, http.StatusInternalServerError)
}
