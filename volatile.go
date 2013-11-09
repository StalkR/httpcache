package httpcache

import (
	"net/http"
)

// NewVolatile creates an http RoundTripper with a memory cache.
func NewVolatile(transport http.RoundTripper) http.RoundTripper {
	return &CachedRoundTrip{
		Transport: transport,
		Cache:     memoryCache{Map: make(map[string][]byte)},
	}
}

// NewVolatileClient creates an http client with a memory cache.
func NewVolatileClient() *http.Client {
	return &http.Client{Transport: NewVolatile(http.DefaultTransport)}
}
