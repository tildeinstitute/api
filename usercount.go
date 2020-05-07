package main

import "net/http"

// UserCount handles the /<format>/usercount endpoint.
// Responds with the number of registered users on the system.
func UserCount(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
