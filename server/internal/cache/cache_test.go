package cache

import (
	"testing"
	"time"
)

func newTestCache(t *testing.T) *Cache {
	t.Helper()
	c, err := New(1e4, 1<<20)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

func TestCache_SetAndGet(t *testing.T) {
	c := newTestCache(t)
	c.Set("key1", "value1", 1)
	// Ristretto is async; wait for value to be processed
	time.Sleep(10 * time.Millisecond)

	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("Get() should find key1")
	}
	if val != "value1" {
		t.Errorf("Get() = %v, want 'value1'", val)
	}
}

func TestCache_GetMissing(t *testing.T) {
	c := newTestCache(t)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for missing key")
	}
}

func TestCache_Del(t *testing.T) {
	c := newTestCache(t)
	c.Set("key2", "value2", 1)
	time.Sleep(10 * time.Millisecond)

	c.Del("key2")
	time.Sleep(10 * time.Millisecond)

	_, ok := c.Get("key2")
	if ok {
		t.Error("Get() should return false after Del()")
	}
}

func TestCache_SetWithTTL(t *testing.T) {
	c := newTestCache(t)
	c.SetWithTTL("ttl-key", "ttl-value", 1, 50*time.Millisecond)
	time.Sleep(10 * time.Millisecond)

	val, ok := c.Get("ttl-key")
	if !ok {
		t.Fatal("Get() should find ttl-key before expiry")
	}
	if val != "ttl-value" {
		t.Errorf("Get() = %v, want 'ttl-value'", val)
	}

	time.Sleep(100 * time.Millisecond)
	_, ok = c.Get("ttl-key")
	if ok {
		t.Error("Get() should return false after TTL expiry")
	}
}
