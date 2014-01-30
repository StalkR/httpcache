package httpcache

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Reporter is an interface that allows reporting of cache performance,
// e.g. send events to a logger or statsd, etc.
// The handlers accept the request as is, so they can report for instance, hits/misses per domain
type Reporter interface {
	OnCacheHit(*http.Request)
	OnCacheMiss(*http.Request)
	OnUncachable(*http.Request)
}

// the default reporter does nothing
type NopReporter struct{}

func (*NopReporter) OnCacheHit(*http.Request)   {}
func (*NopReporter) OnCacheMiss(*http.Request)  {}
func (*NopReporter) OnUncachable(*http.Request) {}

var currentReporter Reporter = &NopReporter{}

// Set the stats reporter of the cache
func SetReporter(r Reporter) {
	currentReporter = r
}

// Cache represents the ability to cache response data from URL.
type Cache interface {
	Get(u *url.URL) (*entry, error)
	Put(u *url.URL, e *entry) error
}

type entry struct {
	Data     []byte
	SaveTime time.Time
	TTL      time.Duration
}

// A CachedRoundTrip implements net/http RoundTripper with a cache.
type CachedRoundTrip struct {
	m         sync.Mutex
	Transport http.RoundTripper
	Cache     Cache
	Policy    CachePolicyProvider
}

// RoundTrip loads from cache if possible or RoundTrips and saves it.
func (c *CachedRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {

	if !c.cacheableRequest(req) {
		//notify the reporter
		currentReporter.OnUncachable(req)

		return c.Transport.RoundTrip(req)
	}
	cache, err := c.load(req, 20)
	if err == nil {
		currentReporter.OnCacheHit(req)
		return cache, nil
	}
	currentReporter.OnCacheMiss(req)

	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if c.cacheableResponse(resp) {
		c.save(req, resp)
	}
	return resp, nil
}

// cacheableRequest tells whether a request is cacheable or not.
func (c CachedRoundTrip) cacheableRequest(req *http.Request) bool {
	return req.Method == "GET"
}

// cacheableResponse tells whether a response is cacheable or not.
func (c CachedRoundTrip) cacheableResponse(resp *http.Response) bool {
	return resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusMovedPermanently
}

// load prepares the response of a request by loading its body from cache.
func (c CachedRoundTrip) load(req *http.Request, maxRedirects int) (*http.Response, error) {
	if maxRedirects == 0 {
		return nil, errors.New("httpcache: Load: max redirects hit")
	}

	entry, err := c.Cache.Get(req.URL)
	if err != nil || entry == nil {

		return nil, err
	}

	// entry expired!
	if entry.TTL >= 0 && entry.SaveTime.Add(entry.TTL).Before(time.Now()) {
		return nil, errors.New("httpcache: TTL expired")
	}

	body := entry.Data

	if strings.HasPrefix(string(body), "REDIRECT:") {
		u, err := url.Parse(strings.TrimPrefix(string(body), "REDIRECT:"))
		if err != nil {
			return nil, err
		}
		req.URL = u
		return c.load(req, maxRedirects-1)
	}

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}, nil
}

// Generate a new cache entry for a given request, calling the policy provider for TTL
func (c *CachedRoundTrip) newEntry(data []byte, resp *http.Response) *entry {
	return &entry{
		data,
		time.Now(),
		c.Policy.GetTTL(resp),
	}
}

// save saves the body of a response corresponding to a request.
func (c *CachedRoundTrip) save(req *http.Request, resp *http.Response) error {
	if resp.StatusCode == http.StatusMovedPermanently {
		u, err := resp.Location()
		if err != nil {
			return err
		}

		err = c.Cache.Put(req.URL, c.newEntry([]byte("REDIRECT:"+u.String()), resp))
		if err != nil {
			return err
		}
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	err = c.Cache.Put(req.URL, c.newEntry(body, resp))
	if err != nil {
		return err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	return nil
}
