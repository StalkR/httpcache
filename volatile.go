package httpcache

import (
	"net/http"
	"time"
)

// NewVolatile creates an http RoundTripper with a memory cache.
func NewVolatile(transport http.RoundTripper, TTL time.Duration, maxItems int) http.RoundTripper {
	return &CachedRoundTrip{
		Transport: transport,
		Cache:     NewMemoryCache(maxItems),
		TTL:       TTL,
	}
}

// NewVolatileClient creates an http client with a memory cache.
func NewVolatileClient(TTL time.Duration, maxItems int) *http.Client {
	return &http.Client{Transport: NewVolatile(http.DefaultTransport, TTL, maxItems)}
}
