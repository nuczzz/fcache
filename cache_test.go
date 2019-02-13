package fcache

import (
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	cache := newMemCache(10, false, 0)
	cache.Set("key1", []byte("123456789"))
	cache.Set("key2", []byte("0"))
	t.Log(cache.Get("key1"))
	t.Log(cache.Get("key2"))
}

func TestDiskCache(t *testing.T) {
	cache := newDiskCache(10, false, "./cache", 0)
	cache.Set("key1", []byte("123456789"))
	cache.Set("key2", []byte("0"))
	t.Log(cache.Get("key1"))
	t.Log(cache.Get("key2"))
}

func TestClear(t *testing.T) {
	memCache := newMemCache(10, false, 0)
	memCache.Set("key1", []byte("123456789"))
	t.Log(memCache.Get("key1"))
	t.Log(memCache.Clear("key1"))
	t.Log(memCache.Get("key1"))

	diskCache := newDiskCache(10, false, "./cache", 0)
	diskCache.Set("key1", []byte("123456789"))
	t.Log(diskCache.Get("key1"))
	time.Sleep(time.Second)
	t.Log(diskCache.Clear("key1"))
	t.Log(diskCache.Get("key1"))
}

func TestClearAll(t *testing.T) {
	memCache := newMemCache(100, false, 0)
	memCache.Set("key1", []byte("123456789"))
	memCache.Set("key2", []byte("123456789"))
	t.Log(memCache.ClearAll())
	t.Log(memCache.Get("key1"))

	diskCache := newDiskCache(100, false, "./cache", 0)
	t.Log(diskCache.Set("key1", []byte("123456789")))
	t.Log(diskCache.Set("key2", []byte("123456789")))
	time.Sleep(time.Second)
	t.Log(diskCache.ClearAll())
	t.Log(diskCache.Get("key1"))
}
