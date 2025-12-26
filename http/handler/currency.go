package handler

import (
	"errors"
	"log"
	"strings"

	"github.com/BohdanKyryliuk/golang/currency_converter"
	"github.com/BohdanKyryliuk/golang/currencyapi"
	"github.com/gin-gonic/gin"
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
func (h *Currency) Status(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	status, err := h.client.CheckStatus(c.Request.Context())
	if err != nil {
		handleCurrencyError(c, err)
		return
	}

	c.String(200, "%s", status)
}

// Currencies handles requests for available currencies
func (h *Currency) Currencies(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	currencies, err := h.client.GetCurrencies(c.Request.Context())
	if err != nil {
		handleCurrencyError(c, err)
		return
	}

	c.String(200, "%s", currencies)
}

// LatestRates handles requests for latest exchange rates
// Query params: base (base currency), currencies (comma-separated list)
func (h *Currency) LatestRates(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	// Parse query parameters
	params := &currency_converter.LatestRatesParams{
		BaseCurrency: c.Query("base"),
	}
	if currencies := c.Query("currencies"); currencies != "" {
		params.Currencies = strings.Split(currencies, ",")
	}

	rates, err := h.client.GetLatestRates(c.Request.Context(), params)
	if err != nil {
		handleCurrencyError(c, err)
		return
	}

	c.String(200, "%s", rates)
}

// handleCurrencyError handles errors from the currency converter with appropriate HTTP responses
func handleCurrencyError(c *gin.Context, err error) {
	log.Printf("Currency API error: %v", err)

	// Check for specific error types and set appropriate status codes
	var apiErr *currencyapi.APIError
	if errors.As(err, &apiErr) {
		if apiErr.IsInvalidAPIKey() {
			c.AbortWithStatusJSON(500, gin.H{"error": "Service configuration error"})
			return
		}
		if apiErr.IsQuotaExceeded() {
			c.AbortWithStatusJSON(503, gin.H{"error": "Service temporarily unavailable, please try again later"})
			return
		}
	}

	var httpErr *currencyapi.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.IsRateLimited() {
			c.AbortWithStatusJSON(429, gin.H{"error": "Rate limited, please try again later"})
			return
		}
	}

	// Check if it's a temporary error
	if currencyapi.IsTemporaryError(err) {
		c.AbortWithStatusJSON(503, gin.H{"error": "Service temporarily unavailable"})
		return
	}

	// Default error response
	c.AbortWithStatusJSON(500, gin.H{"error": "Failed to fetch currency data"})
}
