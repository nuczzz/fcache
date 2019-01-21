package fcache

import (
	"fmt"
	"testing"
	"lru"
)

func TestMemCacheMaxSize(t *testing.T) {
	cache1 := &memCache{
		m:       make(map[interface{}]*lru.Node),
	}
	link := &lru.LRU{
		MaxSize:            100,
		DeleteNodeCallBack: cache1.deleteCallBack(),
	}
	cache1.lru = link
	t.Logf("%#v", cache1)
	// 50 bytes
	for i := 0; i < 9; i++ {
		cache1.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	cache1.Set("key10", []byte("12345"))
	cache1.Set("key11", []byte("123"))
	cache1.Set("key12", []byte("12"))
	t.Logf("%#v", cache1)
	t.Log(len(cache1.m))
	if len(cache1.m) != 12 {
		t.Fatal("maxSize error")
	}
	t.Log(cache1.lru.Traversal())

	cache2 := &memCache{
		m:       make(map[interface{}]*lru.Node),
	}
	link = &lru.LRU{
		MaxSize:            100,
		DeleteNodeCallBack: cache2.deleteCallBack(),
	}
	cache2.lru = link
	t.Logf("%#v", cache2)
	// 50 bytes
	for i := 0; i < 9; i++ {
		cache2.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	cache2.Set("key10", []byte("12345"))
	cache2.Set("key11", []byte("123"))
	cache2.Set("key12", []byte("12"))
	cache2.Set("key13", []byte("1"))
	t.Logf("%#v", cache2)
	t.Log(len(cache2.m))
	if len(cache2.m) != 12 {
		t.Fatal("maxSize error")
	}
	t.Log(cache2.lru.Traversal())
}

func TestMemCache(t *testing.T) {
	cache := &memCache{
		m:       make(map[interface{}]*lru.Node),
		maxSize: 10,
	}
	link := &lru.LRU{
		MaxSize:            100,
		DeleteNodeCallBack: cache.deleteCallBack(),
	}
	cache.lru = link
	t.Logf("%#v", cache)

	cache.Set("key1", []byte("123456789"))
	cache.Set("key2", []byte("0"))
	cache.Set("key3", []byte("1"))
	t.Log(cache.Get("key3"))
	t.Log(cache.Get("key1"))

	t.Logf("%#v", cache)
	t.Log(cache.lru.Traversal())
}
