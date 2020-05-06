package main

import (
	"errors"
	"net/http"
	"strings"
)

const mimePlain = "text/plain; charset=utf-8"
const mimeJSON = "application/json; charset=utf-8"

// Validates the request and then sends it off to where it needs to go.
// Eventually.
func validateRequest(w http.ResponseWriter, r *http.Request) {
	if !methodHop(r) {
		errHTTP(w, r, errors.New("405 Method Not Allowed"), http.StatusMethodNotAllowed)
		return
	}

	format := formatHop(r)
	if format != "plain" && format != "json" {
		errHTTP(w, r, errors.New("400 Bad Request"), http.StatusBadRequest)
		return
	}
}

// Simple HTTP method check
func methodHop(r *http.Request) bool {
	return r.Method == http.MethodGet || r.Method == http.MethodHead
}

// Makes sure the format requested is either plaintext or JSON
func formatHop(r *http.Request) string {
	split := strings.Split(r.URL.Path[1:], "/")
	return split[0]
}

// Chooses next hop based on the endpoint
func routingHop(w http.ResponseWriter, r *http.Request, format string) {

}
