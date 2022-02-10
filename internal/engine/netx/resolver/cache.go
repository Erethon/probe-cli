package resolver

import (
	"context"
	"sync"

	"github.com/ooni/probe-cli/v3/internal/model"
)

// CacheResolver is a resolver that caches successful replies.
type CacheResolver struct {
	Cache    map[string][]string
	ReadOnly bool
	model.Resolver
	mu sync.Mutex
}

// LookupHost implements Resolver.LookupHost
func (r *CacheResolver) LookupHost(
	ctx context.Context, hostname string) ([]string, error) {
	if entry := r.Get(hostname); entry != nil {
		return entry, nil
	}
	entry, err := r.Resolver.LookupHost(ctx, hostname)
	if err != nil {
		return nil, err
	}
	if !r.ReadOnly {
		r.Set(hostname, entry)
	}
	return entry, nil
}

// Get gets the currently configured entry for domain, or nil
func (r *CacheResolver) Get(domain string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Cache[domain]
}

// Set allows to pre-populate the cache
func (r *CacheResolver) Set(domain string, addresses []string) {
	r.mu.Lock()
	if r.Cache == nil {
		r.Cache = make(map[string][]string)
	}
	r.Cache[domain] = addresses
	r.mu.Unlock()
}
