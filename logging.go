package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func initLogging(logChan <-chan struct{}) {
	logfile, err := os.OpenFile("api.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Printf("Log init error: %s", err.Error())
		return
	}

	go func(*os.File) {
		<-logChan
		log.Printf("Closing log file ...")

		if err := logfile.Close(); err != nil {
			log.SetOutput(os.Stderr)
			log.Printf("Error closing log file: %s", err.Error())
		}
	}(logfile)

	log.SetOutput(logfile)
}

// Appends a 200 OK to the request log
func log200(r *http.Request) {
	useragent := r.Header["User-Agent"]
	uip := getIPFromCtx(r.Context())
	log.Printf("*** %v :: 200 :: %v %v :: %v\n", uip, r.Method, r.URL, useragent)
}

// Appends a request of a given status code to the request
// log. Intended for errors.
func errHTTP(w http.ResponseWriter, r *http.Request, err error, code int) {
	useragent := r.Header["User-Agent"]
	uip := getIPFromCtx(r.Context())
	log.Printf("*** %v :: %v :: %v %v :: %v :: %v\n", uip, code, r.Method, r.URL, useragent, err.Error())
	http.Error(w, fmt.Sprintf("Error %v: %v", code, err.Error()), code)
}
