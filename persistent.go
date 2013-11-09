package httpcache

import (
	"net/http"
)

// NewPersistent creates an http RoundTripper with a file cache.
// Cache files will be created under path.
func NewPersistent(transport http.RoundTripper, path string) http.RoundTripper {
	return &CachedRoundTrip{
		Transport: transport,
		Cache:     fileCache{Path: path},
	}
}

// NewPersistentClient creates an http client with a file cache.
func NewPersistentClient(path string) *http.Client {
	return &http.Client{Transport: NewPersistent(http.DefaultTransport, path)}
}
