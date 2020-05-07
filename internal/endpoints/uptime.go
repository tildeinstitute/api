package endpoints // import git.tilde.institute/tilde/api/internal/endpoints

import "net/http"

// Uptime handles the /<format>/uptime endpoint.
// Sends uptime and load
func Uptime(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
