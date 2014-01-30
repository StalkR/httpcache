package httpcache

import (
	"fmt"
	"net/http"
	"os"
)

// NewPersistent creates an http RoundTripper with a file cache.
// Cache files will be created under path.
func NewPersistent(transport http.RoundTripper, path string, policy CachePolicyProvider) http.RoundTripper {
	return &CachedRoundTrip{
		Transport: transport,
		Cache:     fileCache{Path: path},
		Policy:    policy,
	}
}

// NewPersistentClient creates an http client with a file cache.
func NewPersistentClient(path string, policy CachePolicyProvider) *http.Client {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Printf("Could not create %s: %s", path, err)
			return nil
		}
	}

	return &http.Client{Transport: NewPersistent(http.DefaultTransport, path, policy)}
}
