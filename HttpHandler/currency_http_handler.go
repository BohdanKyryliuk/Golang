package HttpHandler

import (
	"Golang/currency_converter"
	"fmt"
	"net/http"
)

func CurrencyStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status := currency_converter.CheckStatus()

	fmt.Fprintf(w, status)
}
