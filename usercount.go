package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Just returns the number of directories in /home
// The assumption being, it's the number of human users.
func userCountQuery(format string) ([]byte, error) {
	ls, err := exec.Command("/bin/ls", "/home").Output()
	if err != nil {
		return nil, fmt.Errorf("Couldn't execute ls: %w", err)
	}

	split := strings.Fields(string(ls))

	if format == "plain" {
		return []byte(fmt.Sprintf("%v users\n", len(split))), nil
	}

	out := fmt.Sprintf(`{
	"userCount": "%v"
}
`, len(split))

	return []byte(out), nil
}
