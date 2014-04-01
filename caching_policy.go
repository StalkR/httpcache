package httpcache

import (
	//	"fmt"
	"net/http"
	"time"
)

// This interface is responsible for telling the client for how long we cache a response.
// For each response the provider returns a TTL as time.Duration, telling the cache for
// how long the response to this request needs to be cached.
// Return < 0 for inifinite caching, or 0 for no caching
//
// By default the cache can
type CachePolicyProvider interface {
	GetTTL(*http.Response) time.Duration
}

const DONT_CACHE time.Duration = 0
const CACHE_FOREVER time.Duration = -1

// NeverExpirePolicy is a policy provider that caches everything forever
type NeverExpirePolicy struct{}

func (NeverExpirePolicy) GetTTL(resp *http.Response) time.Duration {
	return CACHE_FOREVER
}

// A caching policy provider that provides a fixed expiration for all requests
type FixedTTLPolicy struct {
	TTL time.Duration
}

// Create a new fixed ttl policy provider with a given TTL
func NewFixedTTLPolicy(ttl time.Duration) *FixedTTLPolicy {
	return &FixedTTLPolicy{ttl}
}

func (p *FixedTTLPolicy) GetTTL(resp *http.Response) time.Duration {
	return p.TTL
}

// PerDomainTTLPolicy is a caching policy with different expiration per domain.
// It allows you to modify per domain policies in run time
type PerDomainTTLPolicy struct {
	domains    map[string]time.Duration
	defaultTTL time.Duration
}

// Create a new domain based caching policy, given a map of domain=>TTL.
//

func NewPerDomainTTLPolicy(rules map[string]time.Duration, defaultTTL time.Duration) *PerDomainTTLPolicy {
	if rules == nil {
		rules = make(map[string]time.Duration)
	}

	// FIXME: Right now this should be done separately for all sub domains of a domain,
	// but we should do proper normalizing of subdomains
	return &PerDomainTTLPolicy{
		domains:    rules,
		defaultTTL: defaultTTL,
	}
}

// Set TTL for a specific domain
func (p *PerDomainTTLPolicy) SetTTL(domain string, ttl time.Duration) {
	p.domains[domain] = ttl
}

// Return the domain's policy TTL, or the default TTL
func (p *PerDomainTTLPolicy) GetTTL(resp *http.Response) time.Duration {

	domain := resp.Request.URL.Host
	ret, found := p.domains[domain]
	if !found {
		ret = p.defaultTTL
	}

	return ret
}
