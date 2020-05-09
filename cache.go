package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

const cacheDuration = 1 * time.Minute

// Holds the cached responses
type page struct {
	raw     []byte
	expires time.Time
}

// Wraps the page cache map with a rwlock
type cacheWrapper struct {
	sync.RWMutex
	pages map[string]*page
}

// The actual page/response cache
var cache = &cacheWrapper{
	pages: make(map[string]*page),
}

// Checks if a page exists in the cache already.
// If it doesn't, creates an empty entry.
func unNullPage(path string) {
	cache.RLock()
	pageBlob := cache.pages[path]
	cache.RUnlock()

	if pageBlob == nil {
		cache.Lock()
		cache.pages[path] = &page{
			raw:     []byte{},
			expires: time.Time{},
		}
		cache.Unlock()
	}
}

// Returns true if the cached page is good.
// False if it's stale.
func (cache *cacheWrapper) isFresh(path string) bool {
	cache.RLock()
	defer cache.RUnlock()
	return time.Now().Before(cache.pages[path].expires)
}

// Wraps the two cache-checking functions.
// One for /, the other for various requests.
func (cache *cacheWrapper) bap(requestPath string) {
	if requestPath == "/" {
		bapIndex()
		return
	}
	split := strings.Split(requestPath[1:], "/")
	switch split[1] {
	case "osversion":
		bapOSVersion(split[0])
	case "pkgs":
		bapPkgs(split[0])
	case "uptime":
		bapUptime(split[0])
	case "usercount":
		bapUserCount(split[0])
	case "users":
		bapUsers(split[0])
	default:
	}
}

// yoinks the raw data to send to the requester
func (cache *cacheWrapper) yoink(path string) []byte {
	cache.RLock()
	defer cache.RUnlock()

	return cache.pages[path].raw
}

// yoinks the expiration for the cache
func (cache *cacheWrapper) expiresWhen(path string) string {
	cache.RLock()
	defer cache.RUnlock()

	return cache.pages[path].expires.Format(time.RFC1123)
}

// Checks if cache either has expired or has nil copy of
// the index. If so, it yoinks the page from disk and
// sets the expiration time.
func bapIndex() {
	unNullPage("/")

	if cache.isFresh("/") {
		return
	}

	bytes, err := ioutil.ReadFile("web/index.txt")
	if err != nil {
		log.Printf("Could not read index page: %s", err.Error())
		bytes = []byte("tilde.institute informational API")
	}

	cache.Lock()
	defer cache.Unlock()

	cache.pages["/"] = &page{
		raw:     bytes,
		expires: time.Now().Add(cacheDuration),
	}
}

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

func bapPkgs(format string)      {}
func bapQuery(format string)     {}
func bapUptime(format string)    {}
func bapUserCount(format string) {}
func bapUsers(format string)     {}
