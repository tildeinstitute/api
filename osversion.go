package main

import "net/http"

// OSVersion handles the /<format>/osversion endpoint.
// Responds with the OpenBSD version.
func OSVersion(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
