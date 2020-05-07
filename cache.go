package main

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

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

// Checks if cache either has expired or has nil copy of
// the index. If so, it yoinks the page from disk and
// sets the expiration time.
func bapIndex() {
	if cache.pages["/"] == nil {
		cache.Lock()
		cache.pages["/"] = &page{
			raw:     []byte{},
			expires: time.Time{},
		}
		cache.Unlock()
	}

	cache.RLock()
	expires := cache.pages["/"].expires
	cache.RUnlock()

	if time.Now().Before(expires) {
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
		expires: time.Now().Add(1 * time.Minute),
	}
}

func bapOSVersion(format string) {

}

func bapPkgs(format string)      {}
func bapQuery(format string)     {}
func bapUptime(format string)    {}
func bapUserCount(format string) {}
func bapUsers(format string)     {}
