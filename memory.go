package httpcache

import (
	"errors"
	"net/url"
	"time"
)

// A memoryCache implements a cache with files saved in Path.
// It is not safe for concurrency.
type memoryCache struct {
	Map map[string]*entry
}

// Get gets data saved for an URL if present in cache.
func (m memoryCache) Get(u *url.URL) (*entry, error) {
	data, ok := m.Map[u.String()]
	if !ok {
		return nil, errors.New("not in cache")
	}
	return data, nil
}

// Put puts data of an URL in cache.
func (m memoryCache) Put(u *url.URL, data []byte) error {
	m.Map[u.String()] = &entry{data, time.Now()}
	return nil
}
