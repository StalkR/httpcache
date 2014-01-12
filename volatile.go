package httpcache

import (
	"net/http"
	"time"
)

// NewVolatile creates an http RoundTripper with a memory cache.
func NewVolatile(transport http.RoundTripper, TTL time.Duration) http.RoundTripper {
	return &CachedRoundTrip{
		Transport: transport,
		Cache:     memoryCache{Map: make(map[string]*entry)},
		TTL:       TTL,
	}
}

// NewVolatileClient creates an http client with a memory cache.
func NewVolatileClient(TTL time.Duration) *http.Client {
	return &http.Client{Transport: NewVolatile(http.DefaultTransport, TTL)}
}
