package handler

import (
	"errors"
	"strings"

	"github.com/BohdanKyryliuk/golang/worker"
	"github.com/gin-gonic/gin"
)

// Rates holds the dependencies for rate-related HTTP handlers
type Rates struct {
	manager *worker.Manager
}

// NewRates creates a new Rates handler with the given worker manager
func NewRates(manager *worker.Manager) *Rates {
	return &Rates{manager: manager}
}

// GetRate handles requests for cached rates of a specific base currency
// Query params: base (base currency, required)
func (h *Rates) GetRate(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	baseCurrency := strings.ToUpper(c.Query("base"))
	if baseCurrency == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "base currency parameter is required"})
		return
	}

	rateData, err := h.manager.GetRates(baseCurrency)
	if err != nil {
		var notFoundErr *worker.NotFoundError
		if errors.As(err, &notFoundErr) {
			c.AbortWithStatusJSON(404, gin.H{"error": "rates not found for currency: " + baseCurrency})
			return
		}
		c.AbortWithStatusJSON(500, gin.H{"error": "failed to get rates"})
		return
	}

	c.JSON(200, rateData)
}

// GetAllRates handles requests for all cached rates
func (h *Rates) GetAllRates(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	allRates := h.manager.GetAllRates()
	c.JSON(200, allRates)
}

// GetWorkerStatus handles requests for worker status
func (h *Rates) GetWorkerStatus(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	status := struct {
		Running    bool     `json:"running"`
		Currencies []string `json:"currencies"`
	}{
		Running:    h.manager.IsRunning(),
		Currencies: h.manager.GetCurrencies(),
	}

	c.JSON(200, status)
}
