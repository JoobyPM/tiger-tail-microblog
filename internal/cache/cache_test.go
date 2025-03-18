package cache

import (
	"bytes"
	"testing"
	"time"
)

func TestNewMemoryCache(t *testing.T) {
	cache := NewMemoryCache()
	
	if cache == nil {
		t.Fatal("NewMemoryCache returned nil")
	}
	
	if cache.items == nil {
		t.Fatal("cache.items is nil")
	}
	
	if len(cache.items) != 0 {
		t.Errorf("cache.items has %d items, want 0", len(cache.items))
	}
}

func TestMemoryCacheGet(t *testing.T) {
	cache := NewMemoryCache()
	
	// Test getting a non-existent key
	_, err := cache.Get("non-existent")
	if err != ErrNotFound {
		t.Errorf("Get non-existent key error = %v, want %v", err, ErrNotFound)
	}
	
	// Test getting an existing key
	expectedValue := []byte("test value")
	cache.items["test-key"] = item{
		value:      expectedValue,
		expiration: 0, // No expiration
	}
	
	value, err := cache.Get("test-key")
	if err != nil {
		t.Errorf("Get existing key error = %v, want nil", err)
	}
	
	if !bytes.Equal(value, expectedValue) {
		t.Errorf("Get value = %v, want %v", value, expectedValue)
	}
	
	// Test getting an expired key
	cache.items["expired-key"] = item{
		value:      []byte("expired value"),
		expiration: time.Now().Add(-1 * time.Hour).UnixNano(), // Expired 1 hour ago
	}
	
	_, err = cache.Get("expired-key")
	if err != ErrNotFound {
		t.Errorf("Get expired key error = %v, want %v", err, ErrNotFound)
	}
	
	// Verify the expired key was removed
	if _, exists := cache.items["expired-key"]; exists {
		t.Error("Expired key was not removed from the cache")
	}
}

func TestMemoryCacheSet(t *testing.T) {
	cache := NewMemoryCache()
	
	// Test setting a key with no expiration
	err := cache.Set("test-key", []byte("test value"), 0)
	if err != nil {
		t.Errorf("Set error = %v, want nil", err)
	}
	
	item, exists := cache.items["test-key"]
	if !exists {
		t.Error("Key was not added to the cache")
	}
	
	if !bytes.Equal(item.value, []byte("test value")) {
		t.Errorf("Item value = %v, want %v", item.value, []byte("test value"))
	}
	
	if item.expiration != 0 {
		t.Errorf("Item expiration = %v, want 0", item.expiration)
	}
	
	// Test setting a key with expiration
	err = cache.Set("expiring-key", []byte("expiring value"), 1*time.Hour)
	if err != nil {
		t.Errorf("Set with expiration error = %v, want nil", err)
	}
	
	item, exists = cache.items["expiring-key"]
	if !exists {
		t.Error("Expiring key was not added to the cache")
	}
	
	if !bytes.Equal(item.value, []byte("expiring value")) {
		t.Errorf("Expiring item value = %v, want %v", item.value, []byte("expiring value"))
	}
	
	// The expiration should be roughly 1 hour in the future
	expectedExpiration := time.Now().Add(1 * time.Hour).UnixNano()
	if item.expiration < expectedExpiration-int64(10*time.Second) || item.expiration > expectedExpiration+int64(10*time.Second) {
		t.Errorf("Item expiration = %v, want roughly %v", item.expiration, expectedExpiration)
	}
	
	// Test overwriting an existing key
	err = cache.Set("test-key", []byte("new value"), 0)
	if err != nil {
		t.Errorf("Set overwrite error = %v, want nil", err)
	}
	
	item, exists = cache.items["test-key"]
	if !exists {
		t.Error("Overwritten key was not in the cache")
	}
	
	if !bytes.Equal(item.value, []byte("new value")) {
		t.Errorf("Overwritten item value = %v, want %v", item.value, []byte("new value"))
	}
}

func TestMemoryCacheDelete(t *testing.T) {
	cache := NewMemoryCache()
	
	// Add a key to delete
	cache.items["test-key"] = item{
		value:      []byte("test value"),
		expiration: 0,
	}
	
	// Test deleting an existing key
	err := cache.Delete("test-key")
	if err != nil {
		t.Errorf("Delete error = %v, want nil", err)
	}
	
	if _, exists := cache.items["test-key"]; exists {
		t.Error("Key was not deleted from the cache")
	}
	
	// Test deleting a non-existent key (should not error)
	err = cache.Delete("non-existent")
	if err != nil {
		t.Errorf("Delete non-existent key error = %v, want nil", err)
	}
}

func TestMemoryCacheClear(t *testing.T) {
	cache := NewMemoryCache()
	
	// Add some items to the cache
	cache.items["key1"] = item{value: []byte("value1"), expiration: 0}
	cache.items["key2"] = item{value: []byte("value2"), expiration: 0}
	
	// Test clearing the cache
	err := cache.Clear()
	if err != nil {
		t.Errorf("Clear error = %v, want nil", err)
	}
	
	if len(cache.items) != 0 {
		t.Errorf("Cache has %d items after clear, want 0", len(cache.items))
	}
}

func TestMemoryCacheClose(t *testing.T) {
	cache := NewMemoryCache()
	
	// Test closing the cache (should be a no-op for memory cache)
	err := cache.Close()
	if err != nil {
		t.Errorf("Close error = %v, want nil", err)
	}
}

func TestMemoryCacheCleanup(t *testing.T) {
	cache := NewMemoryCache()
	
	// Add some items to the cache, some expired and some not
	cache.items["expired1"] = item{
		value:      []byte("expired1"),
		expiration: time.Now().Add(-1 * time.Hour).UnixNano(), // Expired 1 hour ago
	}
	cache.items["expired2"] = item{
		value:      []byte("expired2"),
		expiration: time.Now().Add(-30 * time.Minute).UnixNano(), // Expired 30 minutes ago
	}
	cache.items["valid1"] = item{
		value:      []byte("valid1"),
		expiration: time.Now().Add(1 * time.Hour).UnixNano(), // Expires in 1 hour
	}
	cache.items["valid2"] = item{
		value:      []byte("valid2"),
		expiration: 0, // Never expires
	}
	
	// Run cleanup
	cache.cleanup()
	
	// Check that expired items were removed
	if _, exists := cache.items["expired1"]; exists {
		t.Error("Expired item 1 was not removed during cleanup")
	}
	if _, exists := cache.items["expired2"]; exists {
		t.Error("Expired item 2 was not removed during cleanup")
	}
	
	// Check that valid items were not removed
	if _, exists := cache.items["valid1"]; !exists {
		t.Error("Valid item 1 was incorrectly removed during cleanup")
	}
	if _, exists := cache.items["valid2"]; !exists {
		t.Error("Valid item 2 was incorrectly removed during cleanup")
	}
}

// Test for RedisCache - since it's a placeholder, we'll just test that the methods return "not implemented"
func TestRedisCache(t *testing.T) {
	cache, err := NewRedisCache("localhost", 6379, "", 0)
	if err != nil {
		t.Errorf("NewRedisCache error = %v, want nil", err)
	}
	
	if cache == nil {
		t.Fatal("NewRedisCache returned nil")
	}
	
	// Test Get
	_, err = cache.Get("test-key")
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("Get error = %v, want 'not implemented'", err)
	}
	
	// Test Set
	err = cache.Set("test-key", []byte("test value"), 0)
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("Set error = %v, want 'not implemented'", err)
	}
	
	// Test Delete
	err = cache.Delete("test-key")
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("Delete error = %v, want 'not implemented'", err)
	}
	
	// Test Clear
	err = cache.Clear()
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("Clear error = %v, want 'not implemented'", err)
	}
	
	// Test Close
	err = cache.Close()
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("Close error = %v, want 'not implemented'", err)
	}
}
