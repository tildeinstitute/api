package endpoints // import git.tilde.institute/tilde/api/internal/endpoints

import "net/http"

// Users handles the /<format>/users endpoint.
// Responds with information on the system's users.
func Users(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
