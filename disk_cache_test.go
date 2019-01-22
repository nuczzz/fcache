package fcache

import (
	"fmt"
	"testing"
	"time"
)

func TestDiskCacheFileName(t *testing.T) {
	cache := &diskCache{dir: DefaultDiskCacheDir}
	t.Log(cache.fileName("test"))
}

func TestDiskCacheCreateFile(t *testing.T) {
	cache := newDiskCache(100, false, DefaultDiskCacheDir, 0).(*diskCache)
	if err := cache.createFile("test", []byte("1234567890")); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", cache)
}

func TestDiskCacheSet(t *testing.T) {
	cache := newDiskCache(100, false, DefaultDiskCacheDir, 0).(*diskCache)
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	t.Logf("%#v", cache)
	t.Log(cache.lru.Traversal())
}

func TestDiskCacheGet(t *testing.T) {
	cache := newDiskCache(100, false, DefaultDiskCacheDir, 0).(*diskCache)
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	t.Logf("%#v", cache)
	t.Log(cache.Get("key5"))
	t.Logf("%#v", cache)
	t.Log(cache.lru.Traversal())
}

func TestDiskCacheEliminate(t *testing.T) {
	cache := newDiskCache(100, false, DefaultDiskCacheDir, 0).(*diskCache)
	for i := 0; i < 10; i++ {
		cache.Set(fmt.Sprintf("key%d", i), []byte("1234567890"))
	}
	t.Logf("%#v", cache)
	t.Log(cache.Get("key5"))
	t.Logf("%#v", cache)
	t.Log(cache.lru.Traversal())
	t.Log(cache.Set("key10", []byte("hello")))
	t.Log(cache.lru.Traversal())
	t.Logf("%#v", cache)
}

func TestDiskCacheInit(t *testing.T) {
	cache := newDiskCache(100, false, DefaultDiskCacheDir, 0).(*diskCache)
	if err := cache.init(); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", cache)
	time.Sleep(time.Second)
	cache.Set("key10", []byte("1111"))
	t.Log(cache.lru.Traversal())
}
