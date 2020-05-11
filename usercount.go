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
		return nil, fmt.Errorf("User Count Query: %w", err)
	}

	split := strings.Fields(string(ls))
	total := 0

	for _, e := range split {
		if strings.HasPrefix(e, ".") || strings.HasPrefix(e, "_") {
			continue
		}
		total++
	}

	if format == "plain" {
		return []byte(fmt.Sprintf("%v users\n", len(split))), nil
	}

	out := fmt.Sprintf(`{
	"userCount": "%v"
}
`, total)

	return []byte(out), nil
}
