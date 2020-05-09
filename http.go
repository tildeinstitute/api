package main

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

const mimePlain = "text/plain; charset=utf-8"
const mimeJSON = "application/json; charset=utf-8"

// Validates the request and then sends it off to where it needs to go.
// Eventually.
// I chose this monolithic handler that calls validation functions to
// determine what to do next because this will make it easier to test.
func mainHandler(w http.ResponseWriter, r *http.Request) {
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

	cache.bap(r.URL.Path)
	out := cache.yoink(r.URL.Path)

	if format == "json" {
		w.Header().Set("Content-Type", mimeJSON)
	} else {
		w.Header().Set("Content-Type", mimePlain)
	}

	w.Header().Set("Expires", cache.expiresWhen(r.URL.Path))

	_, err := w.Write(out)
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
