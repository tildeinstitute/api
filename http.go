package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"git.tilde.institute/tilde/api/internal/endpoints"
)

const mimePlain = "text/plain; charset=utf-8"
const mimeJSON = "application/json; charset=utf-8"

// Validates the request and then sends it off to where it needs to go.
// Eventually.
// I chose this monolithic handler that calls validation functions to
// determine what to do next because this will make it easier to test.
func validateRequest(w http.ResponseWriter, r *http.Request) {
	if !methodHop(r) {
		errHTTP(w, r, errors.New("405 Method Not Allowed"), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		indexHandler(w, r)
		return
	}

	format := formatHop(r)
	if format != "plain" && format != "json" {
		errHTTP(w, r, errors.New("400 Bad Request"), http.StatusBadRequest)
		return
	}

	var err error
	switch routingHop(r) {
	case "pkgs":
		err = endpoints.Pkgs(w, r, format)
	case "query":
		err = endpoints.Query(w, r, format)
	case "uptime":
		err = endpoints.Uptime(w, r, format)
	case "usercount":
		err = endpoints.UserCount(w, r, format)
	case "users":
		err = endpoints.Users(w, r, format)
	case "osversion":
		err = endpoints.OSVersion(w, r, format)
	default:
		errHTTP(w, r, errors.New("Unknown endpoint"), http.StatusNotFound)
		return
	}

	if err != nil {
		errHTTP(w, r, err, http.StatusBadRequest)
		return
	}
	log200(r)
}

// Simple HTTP method check
func methodHop(r *http.Request) bool {
	return r.Method == http.MethodGet || r.Method == http.MethodHead
}

// Yoinks the response format
func formatHop(r *http.Request) string {
	split := strings.Split(r.URL.Path[1:], "/")
	return split[0]
}

// Yoinks the endpoint
func routingHop(r *http.Request) string {
	split := strings.Split(r.URL.Path[1:], "/")
	return split[1]
}

// Yeets the index/summary page to the user
func indexHandler(w http.ResponseWriter, r *http.Request) {
	cache.bap("/")
	cache.RLock()
	defer cache.RUnlock()

	out := cache.pages["/"].raw
	expires := cache.pages["/"].expires

	w.Header().Set("Content-Type", mimePlain)
	w.Header().Set("Expires", expires.Format(time.RFC1123))

	_, err := w.Write(out)
	if err != nil {
		errHTTP(w, r, err, http.StatusInternalServerError)
		return
	}
	log200(r)
}
