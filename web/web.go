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

	log.Println("Listening on :3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
