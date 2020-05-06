package main

import (
	"errors"
	"net/http"
	"strings"
)

const mimePlain = "text/plain; charset=utf-8"
const mimeJSON = "application/json; charset=utf-8"

// Simple HTTP method check
func methodHop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		errHTTP(w, r, errors.New("405 Method Not Allowed"), http.StatusMethodNotAllowed)
		return
	}

	formatHop(w, r)
}

// Makes sure the format requested is either plaintext or JSON
func formatHop(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path[1:], "/")

	if split[0] != "plain" && split[0] != "json" {
		errHTTP(w, r, errors.New("400 Bad Request"), http.StatusBadRequest)
		return
	}

	routingHop(w, r, split[0])
}

// Chooses next hop based on the endpoint
func routingHop(w http.ResponseWriter, r *http.Request, format string) {

}
