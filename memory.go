package httpcache

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

// A memoryCache implements a cache with files saved in Path.
type memoryCache struct {
	Cache *lru.Cache
	mutex sync.Mutex
}

// newMemoryCache creates a new memory cache using groupcache's lru.
func newMemoryCache(maxItems int) *memoryCache {
	return &memoryCache{Cache: lru.New(maxItems)}
}

// Get gets data saved for an URL if present in cache.
func (m *memoryCache) Get(u *url.URL) (*entry, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	data, ok := m.Cache.Get(u.String())
	if !ok {
		return nil, fmt.Errorf("httpcache: %s not in cache", u)
	}
	return data.(*entry), nil
}

// Put puts data of an URL in cache.
func (m *memoryCache) Put(u *url.URL, data []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Cache.Add(u.String(), &entry{data, time.Now()})
	return nil
}
