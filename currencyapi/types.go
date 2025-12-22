package currencyapi

// StatusResponse represents the response from the status endpoint
type StatusResponse struct {
	AccountID int64 `json:"account_id"`
	Quotas    struct {
		Month struct {
			Total     int `json:"total"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
		} `json:"month"`
		Grace struct {
			Total     int `json:"total"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
		} `json:"grace"`
	} `json:"quotas"`
}

// CurrenciesParams represents parameters for the currencies endpoint
type CurrenciesParams struct {
	Currencies []string // List of currency codes to filter
	Type       string   // Currency type filter: "fiat", "crypto", or "metal"
}

// CurrencyInfo represents information about a currency
type CurrencyInfo struct {
	Symbol        string   `json:"symbol"`
	Name          string   `json:"name"`
	SymbolNative  string   `json:"symbol_native"`
	DecimalDigits int      `json:"decimal_digits"`
	Rounding      int      `json:"rounding"`
	Code          string   `json:"code"`
	NamePlural    string   `json:"name_plural"`
	Type          string   `json:"type"`
	Countries     []string `json:"countries"`
}

// CurrenciesResponse represents the response from the currencies endpoint
type CurrenciesResponse struct {
	Data map[string]CurrencyInfo `json:"data"`
}

// LatestParams represents parameters for the latest endpoint
type LatestParams struct {
	BaseCurrency string   // Base currency code (default: USD)
	Currencies   []string // List of currency codes to filter
}

// RateInfo represents exchange rate information
type RateInfo struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

// LatestResponse represents the response from the latest endpoint
type LatestResponse struct {
	Meta struct {
		LastUpdatedAt string `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]RateInfo `json:"data"`
}

// HistoricalParams represents parameters for the historical endpoint
type HistoricalParams struct {
	Date         string   // Required: Date in format YYYY-MM-DD
	BaseCurrency string   // Base currency code (default: USD)
	Currencies   []string // List of currency codes to filter
}

// HistoricalResponse represents the response from the historical endpoint
type HistoricalResponse struct {
	Meta struct {
		LastUpdatedAt string `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]RateInfo `json:"data"`
}

// ConvertParams represents parameters for the convert endpoint
type ConvertParams struct {
	BaseCurrency string   // Base currency code
	Currencies   []string // Target currency codes
	Value        float64  // Amount to convert
	Date         string   // Optional: Date for historical conversion (YYYY-MM-DD)
}

// ConvertRateInfo represents conversion rate information
type ConvertRateInfo struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

// ConvertResponse represents the response from the convert endpoint
type ConvertResponse struct {
	Meta struct {
		LastUpdatedAt string `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]ConvertRateInfo `json:"data"`
}
