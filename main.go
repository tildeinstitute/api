package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logChan := make(chan struct{})
	initLogging(logChan)

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	go func(chan<- struct{}) {
		for range sigC {
			log.Printf("^C Caught. Shutting down ...")
			logChan <- struct{}{}
			os.Exit(1)
		}
	}(logChan)

	mux := http.NewServeMux()
	mux.HandleFunc("/", validateRequest)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	})

	server := &http.Server{
		Addr:         ":9999",
		Handler:      ipMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting up")
	log.Printf("Listening on %s\n", server.Addr)

	log.Fatalf("%s", server.ListenAndServe().Error())
}
