package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// executes uname and responds with "$OS $VERSION" in []byte{}
func osVersionQuery(format string) ([]byte, error) {
	out, err := exec.Command("/usr/bin/uname", "-a").Output()
	if err != nil {
		return nil, errors.New("Couldn't exec `uname -a`")
	}

	split := strings.Split(string(out), " ")

	if format == "json" {
		return []byte(fmt.Sprintf(`{
	"os": "%s",
	"version": "%s"
}
`, split[0], split[2])), nil
	}

	return []byte(fmt.Sprintf("%s %s", split[0], split[2])), nil
}
