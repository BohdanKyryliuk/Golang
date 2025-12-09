package web

import (
	"Golang/HttpHandler"
	"Golang/currency_converter"
	"log"
	"net/http"
)

func StartServer() {
	// Initialize currency converter client from environment variables
	currencyClient, err := currency_converter.NewFromEnv()
	if err != nil {
		log.Printf("Warning: Currency converter not available: %v", err)
		// Continue without currency endpoints
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", HttpHandler.HelloHandler)
	mux.HandleFunc("/count", HttpHandler.CounterHandler)

	// Only register currency handlers if the client is available
	if currencyClient != nil {
		currencyHandler := HttpHandler.NewCurrencyHandler(currencyClient)
		mux.HandleFunc("/currency/status", currencyHandler.StatusHandler)
		mux.HandleFunc("/currency/currencies", currencyHandler.CurrenciesHandler)
		mux.HandleFunc("/currency/latest", currencyHandler.LatestRatesHandler)
	}

	log.Println("Listening on :3001")

	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Fatal(err)
	}
}

// StartServerWithConfig starts the server with a provided currency converter client
// This is useful for testing or when you want more control over the initialization
func StartServerWithConfig(currencyClient *currency_converter.Client) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HttpHandler.HelloHandler)
	mux.HandleFunc("/count", HttpHandler.CounterHandler)

	if currencyClient != nil {
		currencyHandler := HttpHandler.NewCurrencyHandler(currencyClient)
		mux.HandleFunc("/currency/status", currencyHandler.StatusHandler)
		mux.HandleFunc("/currency/currencies", currencyHandler.CurrenciesHandler)
		mux.HandleFunc("/currency/latest", currencyHandler.LatestRatesHandler)
	}

	log.Println("Listening on :3001")

	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Fatal(err)
	}
}
