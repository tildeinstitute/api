package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Users handles the /<format>/users endpoint.
// Responds with information on the system's users.
func usersQuery(format string) ([]byte, error) {
	ls, err := exec.Command("/bin/ls", "/home").Output()
	if err != nil {
		return nil, fmt.Errorf("Users Query: %w", err)
	}

	users := strings.Fields(string(ls))
	out := `{
	"users": [
`
	for i, e := range users {
		if strings.HasPrefix(e, ".") || strings.HasPrefix(e, "_") {
			continue
		}

		out = fmt.Sprintf("%s\t\t\"%s\"", out, e)

		if i < len(users)-1 {
			out = fmt.Sprintf("%s,\n", out)
		} else {
			out = fmt.Sprintf("%s\n", out)
		}
	}

	out = fmt.Sprintf("%s\t]\n}\n", out)

	return []byte(out), nil
}
