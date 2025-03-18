package cache

import (
	"errors"
	"time"
)

// ErrNotFound is returned when a key is not found in the cache
var ErrNotFound = errors.New("key not found in cache")

// Cache represents a cache interface
type Cache interface {
	// Get retrieves a value from the cache
	Get(key string) ([]byte, error)
	
	// Set stores a value in the cache with an optional expiration
	Set(key string, value []byte, expiration time.Duration) error
	
	// Delete removes a value from the cache
	Delete(key string) error
	
	// Clear removes all values from the cache
	Clear() error
	
	// Close closes the cache connection
	Close() error
}

// MemoryCache is an in-memory implementation of the Cache interface
type MemoryCache struct {
	items map[string]item
}

type item struct {
	value      []byte
	expiration int64
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]item),
	}
	
	// Start a goroutine to clean up expired items
	go cache.cleanupLoop()
	
	return cache
}

// Get retrieves a value from the cache
func (c *MemoryCache) Get(key string) ([]byte, error) {
	item, found := c.items[key]
	if !found {
		return nil, ErrNotFound
	}
	
	// Check if the item has expired
	if item.expiration > 0 && item.expiration < time.Now().UnixNano() {
		delete(c.items, key)
		return nil, ErrNotFound
	}
	
	return item.value, nil
}

// Set stores a value in the cache with an optional expiration
func (c *MemoryCache) Set(key string, value []byte, expiration time.Duration) error {
	var exp int64
	
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}
	
	c.items[key] = item{
		value:      value,
		expiration: exp,
	}
	
	return nil
}

// Delete removes a value from the cache
func (c *MemoryCache) Delete(key string) error {
	delete(c.items, key)
	return nil
}

// Clear removes all values from the cache
func (c *MemoryCache) Clear() error {
	c.items = make(map[string]item)
	return nil
}

// Close closes the cache connection
func (c *MemoryCache) Close() error {
	// Nothing to do for in-memory cache
	return nil
}

// cleanupLoop periodically removes expired items from the cache
func (c *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for {
		<-ticker.C
		c.cleanup()
	}
}

// cleanup removes expired items from the cache
func (c *MemoryCache) cleanup() {
	now := time.Now().UnixNano()
	
	for k, v := range c.items {
		if v.expiration > 0 && v.expiration < now {
			delete(c.items, k)
		}
	}
}

// RedisCache is a placeholder for a Redis implementation of the Cache interface
// In a real application, this would use a Redis client library
type RedisCache struct {
	// Redis client would go here
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(host string, port int, password string, db int) (*RedisCache, error) {
	// In a real application, this would initialize a Redis client
	return &RedisCache{}, nil
}

// Get retrieves a value from the cache
func (c *RedisCache) Get(key string) ([]byte, error) {
	// In a real application, this would use the Redis client to get the value
	return nil, errors.New("not implemented")
}

// Set stores a value in the cache with an optional expiration
func (c *RedisCache) Set(key string, value []byte, expiration time.Duration) error {
	// In a real application, this would use the Redis client to set the value
	return errors.New("not implemented")
}

// Delete removes a value from the cache
func (c *RedisCache) Delete(key string) error {
	// In a real application, this would use the Redis client to delete the value
	return errors.New("not implemented")
}

// Clear removes all values from the cache
func (c *RedisCache) Clear() error {
	// In a real application, this would use the Redis client to clear the cache
	return errors.New("not implemented")
}

// Close closes the cache connection
func (c *RedisCache) Close() error {
	// In a real application, this would close the Redis client
	return errors.New("not implemented")
}
