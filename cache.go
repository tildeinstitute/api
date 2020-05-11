package main

import (
	"errors"
	"io/ioutil"
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
func (cache *cacheWrapper) checkedInit(path string) {
	cache.Lock()
	defer cache.Unlock()

	if cache.pages[path] == nil {
		cache.pages[path] = &page{
			raw:     []byte{},
			expires: time.Time{},
		}
	}
}

// Adds a given page to the cache with path as the key
func (cache *cacheWrapper) insert(path string, raw []byte) {
	cache.Lock()
	defer cache.Unlock()

	cache.pages[path] = &page{
		raw:     raw,
		expires: time.Now().Add(cacheDuration),
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
func (cache *cacheWrapper) bap(requestPath string) error {
	cache.checkedInit(requestPath)

	if cache.isFresh(requestPath) {
		return nil
	}

	var format, query string

	if requestPath != "/" {
		split := strings.Split(requestPath[1:], "/")
		format = split[0]
		query = split[1]
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
	case "users":
		bytes, err = usersQuery(format)
	case "pkgs":
		bytes, err = pkgsQuery(format)
	default:
		if requestPath == "/" {
			bytes, err = ioutil.ReadFile("web/index.txt")
		} else {
			err = errors.New("Invalid Query Type")
		}
	}

	if err != nil {
		return err
	}

	cache.insert(requestPath, bytes)
	return nil
}

// yoinks the raw data and expiry time to send to the requester
func (cache *cacheWrapper) yoink(path string) ([]byte, string) {
	cache.RLock()
	defer cache.RUnlock()

	return cache.pages[path].raw, cache.pages[path].expires.Format(time.RFC1123)
}
