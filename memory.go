package httpcache

import (
	"errors"
	"github.com/golang/groupcache/lru"
	"net/url"
	"time"
)

// A memoryCache implements a cache with files saved in Path.
// It is not safe for concurrency.
type memoryCache struct {
	Cache *lru.Cache
}

func NewMemoryCache(maxItems int) *memoryCache {

	return &memoryCache{
		Cache: lru.New(maxItems),
	}
}

// Get gets data saved for an URL if present in cache.
func (m memoryCache) Get(u *url.URL) (*entry, error) {
	data, ok := m.Cache.Get(u.String())
	if !ok {
		return nil, errors.New("not in cache")
	}
	return data.(*entry), nil
}

// Put puts data of an URL in cache.
func (m memoryCache) Put(u *url.URL, data []byte) error {
	m.Cache.Add(u.String(), &entry{data, time.Now()})
	return nil
}
