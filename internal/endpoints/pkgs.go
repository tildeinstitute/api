package endpoints // import git.tilde.institute/tilde/api/internal/endpoints

import "net/http"

// Pkgs handles the /<format>/pkgs endpoint.
// Sends a list of installed packages.
func Pkgs(w http.ResponseWriter, r *http.Request, format string) error {

	return nil
}
