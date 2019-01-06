package hpcache

import (
	"fmt"
	"testing"
)

func traceDiskCacheLinkedList(dc *diskCache, t *testing.T) {
	temp := dc.header
	for temp != nil {
		t.Logf("%#v", temp)
		temp = temp.next
	}
}

func TestDiskCacheFileName(t *testing.T) {
	cache := &diskCache{dir: defaultDiskCachePath}
	t.Log(cache.fileName("test"))
}

func TestDiskCacheCreateFile(t *testing.T) {
	cache := &diskCache{
		dir:     defaultDiskCachePath,
		m:       make(map[string]*diskData),
		maxSize: 100, //bytes
	}
	if err := cache.createFile("test", []byte("1234567890")); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", cache)
}

func TestDiskCacheSet(t *testing.T) {
	cache := &diskCache{
		dir:     defaultDiskCachePath,
		m:       make(map[string]*diskData),
		maxSize: 100, //bytes
	}
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	t.Logf("%#v", cache)
	traceDiskCacheLinkedList(cache, t)
}

func TestDiskCacheGet(t *testing.T) {
	cache := &diskCache{
		dir:     defaultDiskCachePath,
		m:       make(map[string]*diskData),
		maxSize: 100, //bytes
	}
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	t.Logf("%#v", cache)
	t.Log(string(cache.Get("key5")))
	t.Logf("%#v", cache)
	traceDiskCacheLinkedList(cache, t)
}

func TestDiskCacheEliminate(t *testing.T) {
	cache := &diskCache{
		dir:     defaultDiskCachePath,
		m:       make(map[string]*diskData),
		maxSize: 100, //bytes
	}
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	t.Logf("%#v", cache)
	t.Log(string(cache.Get("key5")))
	t.Logf("%#v", cache)
	traceDiskCacheLinkedList(cache, t)
	cache.Set("key10", []byte("hello"))
	traceDiskCacheLinkedList(cache, t)
	t.Logf("%#v", cache)
}
