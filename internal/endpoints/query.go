package endpoints // import git.tilde.institute/tilde/api/internal/endpoints

import "net/http"

// Query handles the /<format>/query endpoint.
// Accept a query param and responds with the appropriate info.
// 		?pkg=$PACKAGENAME
func Query(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
