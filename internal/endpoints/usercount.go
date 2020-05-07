package endpoints // import git.tilde.institute/tilde/api/internal/endpoints

import "net/http"

// UserCount handles the /<format>/usercount endpoint.
// Responds with the number of registered users on the system.
func UserCount(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
