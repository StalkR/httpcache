package httpcache

import (
	"net/http"
)

// NewVolatile creates an http RoundTripper with a memory cache.
func NewVolatile(transport http.RoundTripper, policy CachePolicyProvider, maxItems int) http.RoundTripper {
	return &CachedRoundTrip{
		Transport: transport,
		Cache:     newMemoryCache(maxItems),
		Policy:    policy,
	}
}

// NewVolatileClient creates an http client with a memory cache.
func NewVolatileClient(policy CachePolicyProvider, maxItems int) *http.Client {
	return &http.Client{Transport: NewVolatile(http.DefaultTransport, policy, maxItems)}
}
