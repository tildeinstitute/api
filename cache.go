package main

import (
	"errors"
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
	format := split[0]
	query := split[1]

	unNullPage(requestPath)

	if cache.isFresh(requestPath) {
		return
	}

	var bytes []byte
	var err error

	switch query {
	case "osversion":
		bytes, err = osVersionQuery(format)
	case "uptime":
		bytes, err = uptimeQuery(format)
	case "usercount":
		bytes, err = userCountQuery(format)
	default:
		err = errors.New("Invalid Query Type")
	}

	if err != nil {
		log.Printf("Could not query %s: %s", requestPath, err.Error())
		bytes = []byte("Internal Error")
	}

	cache.Lock()
	defer cache.Unlock()

	cache.pages[requestPath] = &page{
		raw:     bytes,
		expires: time.Now().Add(cacheDuration),
	}
}

// yoinks the raw data and expiry time to send to the requester
func (cache *cacheWrapper) yoink(path string) ([]byte, string) {
	cache.RLock()
	defer cache.RUnlock()

	return cache.pages[path].raw, cache.pages[path].expires.Format(time.RFC1123)
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
