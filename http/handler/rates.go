package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/BohdanKyryliuk/golang/worker"
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
func (h *Rates) GetRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	baseCurrency := strings.ToUpper(r.URL.Query().Get("base"))
	if baseCurrency == "" {
		http.Error(w, `{"error": "base currency parameter is required"}`, http.StatusBadRequest)
		return
	}

	rateData, err := h.manager.GetRates(baseCurrency)
	if err != nil {
		var notFoundErr *worker.NotFoundError
		if errors.As(err, &notFoundErr) {
			http.Error(w, `{"error": "rates not found for currency: `+baseCurrency+`"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "failed to get rates"}`, http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(rateData)
	if err != nil {
		http.Error(w, `{"error": "failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonBytes)
}

// GetAllRates handles requests for all cached rates
func (h *Rates) GetAllRates(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	allRates := h.manager.GetAllRates()

	jsonBytes, err := json.Marshal(allRates)
	if err != nil {
		http.Error(w, `{"error": "failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonBytes)
}

// GetWorkerStatus handles requests for worker status
func (h *Rates) GetWorkerStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status := struct {
		Running    bool     `json:"running"`
		Currencies []string `json:"currencies"`
	}{
		Running:    h.manager.IsRunning(),
		Currencies: h.manager.GetCurrencies(),
	}

	jsonBytes, err := json.Marshal(status)
	if err != nil {
		http.Error(w, `{"error": "failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonBytes)
}
