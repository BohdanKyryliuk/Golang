package web

import (
	"Golang/HttpHandler"
	"log"
	"net/http"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HttpHandler.HelloHandler)
	mux.HandleFunc("/count", HttpHandler.CounterHandler)
	mux.HandleFunc("/currency/status", HttpHandler.CurrencyStatusHandler)

	log.Println("Listening on :3001")

	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Fatal(err)
	}
}
