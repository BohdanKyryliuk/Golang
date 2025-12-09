package HttpHandler

import (
	"Golang/currency_converter"
	"Golang/currencyapi"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// CurrencyHandler holds the dependencies for currency-related HTTP handlers
type CurrencyHandler struct {
	client *currency_converter.Client
}

// NewCurrencyHandler creates a new CurrencyHandler with the given client
func NewCurrencyHandler(client *currency_converter.Client) *CurrencyHandler {
	return &CurrencyHandler{client: client}
}

// StatusHandler handles requests for currency API status
func (h *CurrencyHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status, err := h.client.CheckStatus(r.Context())
	if err != nil {
		handleCurrencyError(w, err)
		return
	}

	fmt.Fprintf(w, "%s", status)
}

// CurrenciesHandler handles requests for available currencies
func (h *CurrencyHandler) CurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	currencies, err := h.client.GetCurrencies(r.Context())
	if err != nil {
		handleCurrencyError(w, err)
		return
	}

	fmt.Fprintf(w, "%s", currencies)
}

// LatestRatesHandler handles requests for latest exchange rates
func (h *CurrencyHandler) LatestRatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	rates, err := h.client.GetLatestRates(r.Context())
	if err != nil {
		handleCurrencyError(w, err)
		return
	}

	fmt.Fprintf(w, "%s", rates)
}

// CurrencyStatusHandler is a legacy handler for backward compatibility
// Deprecated: Use NewCurrencyHandler().StatusHandler instead
func CurrencyStatusHandler(client *currency_converter.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		status, err := client.CheckStatus(context.Background())
		if err != nil {
			handleCurrencyError(w, err)
			return
		}

		fmt.Fprintf(w, "%s", status)
	}
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
