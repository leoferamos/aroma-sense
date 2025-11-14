package shipping

import (
	"sync"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
)

// QuoteCache abstracts caching for shipping quotes.
type QuoteCache interface {
	Get(key string) ([]dto.ShippingOption, bool)
	Set(key string, val []dto.ShippingOption)
}

// InMemoryQuoteCache is a simple TTL-based in-memory cache implementation.
type InMemoryQuoteCache struct {
	mu  sync.Mutex
	ttl time.Duration
	m   map[string]cacheEntry
}

// cacheEntry wraps a cached item with expiry timestamp.
type cacheEntry struct {
	expiresAt time.Time
	value     []dto.ShippingOption
}

func NewInMemoryQuoteCache(ttl time.Duration) *InMemoryQuoteCache {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	return &InMemoryQuoteCache{
		ttl: ttl,
		m:   make(map[string]cacheEntry),
	}
}

// SetTTL updates the cache TTL.
func (c *InMemoryQuoteCache) SetTTL(ttl time.Duration) {
	if ttl > 0 {
		c.mu.Lock()
		c.ttl = ttl
		c.mu.Unlock()
	}
}

// Get retrieves a value from the cache if present and not expired.
func (c *InMemoryQuoteCache) Get(key string) ([]dto.ShippingOption, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.m == nil {
		return nil, false
	}
	if e, ok := c.m[key]; ok {
		if time.Now().Before(e.expiresAt) {
			return e.value, true
		}
		delete(c.m, key)
	}
	return nil, false
}

// Set stores a value in the cache with the configured TTL.
func (c *InMemoryQuoteCache) Set(key string, val []dto.ShippingOption) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.m == nil {
		c.m = make(map[string]cacheEntry)
	}
	c.m[key] = cacheEntry{expiresAt: time.Now().Add(c.ttl), value: val}
}
