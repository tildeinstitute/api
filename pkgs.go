package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Returns a list of packages installed on the system
func pkgsQuery(format string) ([]byte, error) {
	raw, err := exec.Command("pkg_info", "-a").Output()
	if err != nil {
		return nil, err
	}

	if format == "plain" {
		return raw, nil
	}

	json := `{
	"packages": [`
	rawlines := strings.Split(string(raw), "\n")
	for _, line := range rawlines {
		split := strings.Fields(line)
		if len(split) < 2 {
			continue
		}
		desc := strings.Join(split[1:], " ")
		json = fmt.Sprintf(`%s
		{
			"package": "%s",
			"description": "%s"
		},`, json, split[0], desc)
	}
	json = fmt.Sprintf(`%s
	]
}
`, json[:len(json)-1])
	return []byte(json), nil
}
