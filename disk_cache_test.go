package fcache

import (
	"fmt"
	"testing"
	"time"
)

func traceDiskCacheLinkedList(dc *diskCache, t *testing.T) {
	temp := dc.header
	for temp != nil {
		t.Logf("%#v", temp)
		temp = temp.next
	}
}

func TestDiskCacheFileName(t *testing.T) {
	cache := &diskCache{dir: DefaultDiskCacheDir}
	t.Log(cache.fileName("test"))
}

func TestDiskCacheCreateFile(t *testing.T) {
	cache := &diskCache{
		dir:     DefaultDiskCacheDir,
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
		dir:     DefaultDiskCacheDir,
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
		dir:     DefaultDiskCacheDir,
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
		dir:     DefaultDiskCacheDir,
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

func TestDiskCacheInit(t *testing.T) {
	cache := &diskCache{
		dir:     DefaultDiskCacheDir,
		m:       make(map[string]*diskData),
		maxSize: 50, //bytes
	}
	if err := cache.init(); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", cache)
	time.Sleep(time.Second)
	cache.Set("key10", []byte("1111"))
	traceDiskCacheLinkedList(cache, t)
}
