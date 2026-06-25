package cache

import (
	"context"
	"sync"
	"time"

	"github.com/pactus-project/gopkg/scheduler"
)

// BasicCache is a generic, concurrency-safe in-memory cache backed by sync.Map.
type BasicCache[K any, V any] struct {
	cache sync.Map
}

type basicCacheEntry[V any] struct {
	Value  V
	Expiry time.Time
}

// NewBasic creates a new BasicCache and starts a background cleanup goroutine.
func NewBasic[K any, V any](ctx context.Context, opts ...Option) Cache[K, V] {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	cache := &BasicCache[K, V]{
		cache: sync.Map{},
	}

	scheduler.Every(cfg.cleanUpInterval).Do(ctx, func(context.Context) {
		cache.cleanupExpiredEntries()
	})

	return cache
}

// Add stores a new entry in the cache with the given expiration.
//
//   - expiration: 0 for no expiry.
func (c *BasicCache[K, V]) Add(key K, value V, expiration time.Duration) bool {
	var expiry time.Time
	if expiration != 0 {
		expiry = time.Now().Add(expiration)
	}

	entry := basicCacheEntry[V]{Value: value, Expiry: expiry}
	c.cache.Store(key, entry)

	return true
}

// Get retrieves an entry from the cache by key.
//
func (c *BasicCache[K, V]) Get(key K) (V, bool) {
	var zeroV V // zero Value of type V
	value, ok := c.cache.Load(key)
	if !ok {
		return zeroV, false
	}

	return value.(basicCacheEntry[V]).Value, true
}

// Update modifies the value and optionally the expiration of an existing entry.
func (c *BasicCache[K, V]) Update(key K, newValue V, expiration time.Duration) bool {
	// Check if the key exists in the cache
	value, exists := c.cache.Load(key)
	if !exists {
		return false // Key not found, nothing to update
	}

	// Update the Value
	entry, exists := value.(basicCacheEntry[V])
	if !exists {
		return false
	}
	entry.Value = newValue

	// Update the expiration time if a new expiration is provided
	if expiration != 0 {
		entry.Expiry = time.Now().Add(expiration)
	}

	// Store the updated entry back in the cache
	c.cache.Store(key, entry)

	return true
}

// Exists reports whether the given key is present in the cache.
func (c *BasicCache[K, V]) Exists(key K) bool {
	_, ok := c.cache.Load(key)

	return ok
}

// Keys returns a slice of all keys currently stored in the cache.
func (c *BasicCache[K, V]) Keys() []K {
	keys := make([]K, 0)
	c.cache.Range(func(key, _ any) bool {
		keys = append(keys, key.(K))

		return true
	})

	return keys
}

// Delete removes an entry from the cache by key.
func (c *BasicCache[K, V]) Delete(key K) bool {
	c.cache.Delete(key)
	_, ok := c.cache.Load(key)

	return !ok
}

func (c *BasicCache[K, V]) cleanupExpiredEntries() {
	c.cache.Range(func(key, value any) bool {
		entry, ok := value.(basicCacheEntry[V])
		if !ok {
			return true
		}

		// Skip entries with zero expiration time
		if entry.Expiry.IsZero() {
			return true
		}

		if time.Now().After(entry.Expiry) {
			c.cache.Delete(key)
		}

		return true
	})
}
