package fcache

import "testing"

func TestMemoryCache(t *testing.T) {
	cache := newMemCache(10, false, 0)
	cache.Set("key1", []byte("123456789"))
	cache.Set("key2", []byte("0"))
	t.Log(cache.Get("key1"))
	t.Log(cache.Get("key2"))
}

func TestDiskCache(t *testing.T) {
	// todo
}
