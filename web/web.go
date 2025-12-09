package web

import (
	"log"
	"net/http"

	"Golang/currency_converter"
	"Golang/http/handler"
)

func StartServer() {
	// Initialize currency converter client from environment variables
	currencyClient, err := currency_converter.NewFromEnv()
	if err != nil {
		log.Printf("Warning: Currency converter not available: %v", err)
		// Continue without currency endpoints
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Hello)
	mux.HandleFunc("/count", handler.Counter)

	// Only register currency handlers if the client is available
	if currencyClient != nil {
		currencyHandler := handler.NewCurrency(currencyClient)
		mux.HandleFunc("/currency/status", currencyHandler.Status)
		mux.HandleFunc("/currency/currencies", currencyHandler.Currencies)
		mux.HandleFunc("/currency/latest", currencyHandler.LatestRates)
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
	mux.HandleFunc("/", handler.Hello)
	mux.HandleFunc("/count", handler.Counter)

	if currencyClient != nil {
		currencyHandler := handler.NewCurrency(currencyClient)
		mux.HandleFunc("/currency/status", currencyHandler.Status)
		mux.HandleFunc("/currency/currencies", currencyHandler.Currencies)
		mux.HandleFunc("/currency/latest", currencyHandler.LatestRates)
	}

	log.Println("Listening on :3001")

	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Fatal(err)
	}
}
