package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Checks the cache for the freshness of the osversion query.
// If stale, store the query again.
func bapOSVersion(format string) {
	path := fmt.Sprintf("/%s/osversion", format)
	unNullPage(path)

	if cache.isFresh(path) {
		return
	}

	bytes, err := osVersionQuery(format)
	if err != nil {
		log.Printf("Could not query OS version: %s", err.Error())
		bytes = []byte("Internal Error")
	}

	cache.Lock()
	defer cache.Unlock()

	cache.pages[path] = &page{
		raw:     bytes,
		expires: time.Now().Add(cacheDuration),
	}
}

// executes uname and responds with "$OS $VERSION" in []byte{}
func osVersionQuery(format string) ([]byte, error) {
	out, err := exec.Command("/usr/bin/uname", "-a").Output()
	if err != nil {
		return nil, errors.New("Couldn't exec `uname -a`")
	}

	split := strings.Split(string(out), " ")

	if format == "json" {
		return []byte(fmt.Sprintf(`
{
	"os": "%s",
	"version": "%s"
}`, split[0], split[2])), nil
	}

	return []byte(fmt.Sprintf("%s %s", split[0], split[2])), nil
}
