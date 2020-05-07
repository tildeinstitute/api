package endpoints // import git.tilde.institute/tilde/api/internal/endpoints

import "net/http"

// OSVersion handles the /<format>/osversion endpoint.
// Responds with the OpenBSD version.
func OSVersion(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
