package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	initLogging()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	go func() {
		for range sigC {
			log.Printf("^C Caught. Shutting down ...")
			logChan <- struct{}{}
			os.Exit(1)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", validateRequest)

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
