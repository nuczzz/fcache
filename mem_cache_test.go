package hpcache

import (
	"fmt"
	"testing"
)

func traceLRULinkedList(mc *memCache, t *testing.T) {
	temp := mc.header
	for temp != nil {
		t.Logf("key: %s, value: %s, count: %d", temp.key, string(temp.value), temp.accessCount)
		temp = temp.next
	}
}

func TestMemCacheMaxSize(t *testing.T) {
	cache1 := &memCache{
		m:       make(map[interface{}]*memData),
		maxSize: 100, //100 bytes
	}
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
	traceLRULinkedList(cache1, t)

	cache2 := &memCache{
		m:       make(map[interface{}]*memData),
		maxSize: 100, //100 bytes
	}
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
	traceLRULinkedList(cache2, t)
}

func TestMemCache(t *testing.T) {
	cache := &memCache{
		m:       make(map[interface{}]*memData),
		maxSize: 100,
	}
	t.Logf("%#v", cache)

	for i := 0; i < 9; i++ {
		cache.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	t.Log(len(cache.m))

	cache.Get("key3")
	t.Logf("%#v", cache)

	traceLRULinkedList(cache, t)
}
