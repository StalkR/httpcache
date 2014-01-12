package httpcache

import (
	"bytes"
	"errors"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Cache represents the ability to cache response data from URL.
type Cache interface {
	Get(u *url.URL) (*entry, error)
	Put(u *url.URL, data []byte) error
}

type entry struct {
	Data     []byte
	SaveTime time.Time
}

// A CachedRoundTrip implements net/http RoundTripper with a cache.
type CachedRoundTrip struct {
	m         sync.Mutex
	Transport http.RoundTripper
	Cache     Cache
	TTL       time.Duration
}

// RoundTrip loads from cache if possible or RoundTrips and saves it.
func (c *CachedRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	//c.m.Lock()
	//defer c.m.Unlock()
	if !c.cacheableRequest(req) {
		return c.Transport.RoundTrip(req)
	}
	cache, err := c.load(req, 20)
	if err == nil {
		return cache, nil
	}
	log.Println("Error while fetching from cache: ", err)
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

	//entry expired!
	if entry.SaveTime.Add(c.TTL).Before(time.Now()) {
		return nil, errors.New("Cache TTL Expired!")
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

// save saves the body of a response corresponding to a request.
func (c *CachedRoundTrip) save(req *http.Request, resp *http.Response) error {
	if resp.StatusCode == http.StatusMovedPermanently {
		u, err := resp.Location()
		if err != nil {
			return err
		}
		err = c.Cache.Put(req.URL, []byte("REDIRECT:"+u.String()))
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
	err = c.Cache.Put(req.URL, body)
	if err != nil {
		return err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	return nil
}
