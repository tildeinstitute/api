package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// executes uptime and responds with a []byte{}
// TODO: account for uptimes <1d and >1d
func uptimeQuery(format string) ([]byte, error) {
	out, err := exec.Command("/usr/bin/uptime").Output()
	if err != nil {
		return nil, fmt.Errorf("Could not execute /usr/bin/uptime: %w", err)
	}

	split := strings.Split(string(out), ",")

	if format == "plain" {
		cut := fmt.Sprintf("%s,%s,%s,%s", split[0], split[1], split[3], split[4])
		return []byte(cut), nil
	}

	timeup := strings.Split(string(split[0]), " ")
	time := timeup[0]
	up := fmt.Sprintf("%s,%s", timeup[1:], split[1])

	loadcomb := strings.Split(string(split[3]), " ")
	loads := fmt.Sprintf("%s,%s,%s", loadcomb[1], split[3], split[4])

	json := fmt.Sprintf(`
{
	"time": "%s",
	"up": "%s",
	"load": "%s"
}`, time, up, loads)

	return []byte(json), nil
}
