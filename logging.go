package main

import (
	"fmt"
	"log"
	"net/http"
)

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
