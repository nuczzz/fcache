package hpcache

import (
	"sync"
	"time"
)

// memData memory cache data
type memData struct {
	// value value of memory cache
	value interface{}

	// accessTime timestamp of access time
	accessTime int64

	// accessCount access count
	accessCount int64

	//expireTime  int64
}

// memCache memory cache.
type memCache struct {
	// m map of memory cache. the key is file full path name of request,
	// and map value is file data info.
	m map[interface{}]*memData

	// lock lock of memory cache data.
	lock sync.RWMutex

	// maxSize max size of memory cache data.
	maxSize int
}

// Set set memory cache with key-value pair, and covered if key already exist.
func (mc *memCache) Set(key interface{}, value interface{}) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	// todo： overflow

	mc.m[key] = &memData{
		value:       value,
		accessTime:  time.Now().Unix(),
		accessCount: 1,
	}
}

// Get get memory cache with key, a error will be return if key is not exist.
func (mc *memCache) Get(key interface{}) interface{} {
	mc.lock.RLock()
	defer mc.lock.RUnlock()
	if data, ok := mc.m[key]; ok {
		data.accessTime = time.Now().Unix()
		data.accessCount++
		return data.value
	}
	return nil
}

func newMemCache() Cache {
	return &memCache{m: make(map[interface{}]*memData)}
}
