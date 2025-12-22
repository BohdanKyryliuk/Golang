package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/BohdanKyryliuk/golang/currency_converter"
	"github.com/BohdanKyryliuk/golang/currencyapi"
)

// Currency holds the dependencies for currency-related HTTP handlers
type Currency struct {
	client *currency_converter.Client
}

// NewCurrency creates a new Currency handler with the given client
func NewCurrency(client *currency_converter.Client) *Currency {
	return &Currency{client: client}
}

// Status handles requests for currency API status
func (h *Currency) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status, err := h.client.CheckStatus(r.Context())
	if err != nil {
		handleCurrencyError(w, err)
		return
	}

	fmt.Fprintf(w, "%s", status)
}

// Currencies handles requests for available currencies
func (h *Currency) Currencies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	currencies, err := h.client.GetCurrencies(r.Context())
	if err != nil {
		handleCurrencyError(w, err)
		return
	}

	fmt.Fprintf(w, "%s", currencies)
}

// LatestRates handles requests for latest exchange rates
// Query params: base (base currency), currencies (comma-separated list)
func (h *Currency) LatestRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Parse query parameters
	params := &currency_converter.LatestRatesParams{
		BaseCurrency: r.URL.Query().Get("base"),
	}
	if currencies := r.URL.Query().Get("currencies"); currencies != "" {
		params.Currencies = strings.Split(currencies, ",")
	}

	rates, err := h.client.GetLatestRates(r.Context(), params)
	if err != nil {
		handleCurrencyError(w, err)
		return
	}

	fmt.Fprintf(w, "%s", rates)
}

// CurrencyStatus is a legacy handler for backward compatibility
// Deprecated: Use NewCurrency().Status instead
func CurrencyStatus(client *currency_converter.Client) http.HandlerFunc {
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
