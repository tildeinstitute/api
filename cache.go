package main

import (
	"io/ioutil"
	"log"
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

func bapCache(requestPath string) {
	cache.RLock()
	if cache.pages[requestPath] == nil {
		cache.RUnlock()
		cacheIndex()
		return
	}

	expires := cache.pages[requestPath].expires
	cache.RUnlock()

	if time.Now().After(expires) {
		cacheIndex()
	}
}

// Pulls the index page from disk and places it into the cache.
// etag is an fnv32 hash of the raw file bytes, truncated if necessary.
func cacheIndex() {
	bytes, err := ioutil.ReadFile("web/index.txt")
	if err != nil {
		log.Printf("Could not read index page: %s", err.Error())
		bytes = []byte("tilde.institute informational API")
	}

	cache.Lock()
	cache.pages["/"] = &page{
		raw:     bytes,
		expires: time.Now().Add(5 * time.Minute),
	}
	cache.Unlock()
}
