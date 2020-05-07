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

// Wraps the two cache-checking functions.
// One for /, the other for various requests.
func (cache *cacheWrapper) bap(requestPath string) {
	switch requestPath {
	case "/":
		bapIndex()
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
		expires: time.Now().Add(5 * time.Minute),
	}
}
