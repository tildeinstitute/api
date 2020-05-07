package main

import "net/http"

// Query handles the /<format>/query endpoint.
// Accept a query param and responds with the appropriate info.
// 		?pkg=$PACKAGENAME
func Query(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
