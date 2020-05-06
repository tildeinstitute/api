package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", validateRequest)

	server := &http.Server{
		Addr:         ":9999",
		Handler:      ipMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatalf("%s", server.ListenAndServe())
}
