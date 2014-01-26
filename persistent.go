package httpcache

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// NewPersistent creates an http RoundTripper with a file cache.
// Cache files will be created under path.
func NewPersistent(transport http.RoundTripper, path string, TTL time.Duration) http.RoundTripper {
	return &CachedRoundTrip{
		Transport: transport,
		Cache:     &fileCache{Path: path},
		TTL:       TTL,
	}
}

// NewPersistentClient creates an http client with a file cache.
func NewPersistentClient(path string, TTL time.Duration) (*http.Client, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModeDir|os.ModePerm); err != nil {
			return nil, fmt.Errorf("httpcache: could not create dir %s: %v", path, err)
		}
	}
	return &http.Client{Transport: NewPersistent(http.DefaultTransport, path, TTL)}, nil
}
