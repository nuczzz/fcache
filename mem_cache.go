package fcache

import (
	"sync"

	"github.com/nuczzz/lru"
)

// memCache memory cache.
type memCache struct {
	// m map of memory cache.
	m map[interface{}]*lru.Node

	// needCryptKey crypt key or not.
	needCryptKey bool

	// lru lru control
	lru *lru.LRU

	// lock lock of memory cache data.
	lock sync.RWMutex

	// hitCount hit cache count
	hitCount int64

	// totalCount total count, contains hit count and missing count
	totalCount int64
}

func (mc *memCache) deleteCallBack() func(key interface{}) error {
	return func(key interface{}) error {
		delete(mc.m, key)
		return nil
	}
}

func (mc *memCache) addNodeCallback() func(*lru.Node) {
	return func(node *lru.Node) {
		mc.m[node.Key] = node
	}
}

// Set set memory cache with key-value pair, and covered if key already exist.
func (mc *memCache) Set(key string, value []byte) error {
	if mc.needCryptKey {
		key = MD5(key)
	}

	mc.lock.Lock()
	defer mc.lock.Unlock()

	v := CacheValue{Value: value}
	if data, ok := mc.m[key]; ok {
		return mc.lru.Replace(data, v)
	}
	// memory cache ignore this error
	return mc.lru.AddNewNode(key, v)
}

// Get get memory cache with key, a error will be return if key is not exist.
func (mc *memCache) Get(key string) ([]byte, error) {
	if mc.needCryptKey {
		key = MD5(key)
	}

	mc.lock.RLock()
	defer mc.lock.RUnlock()

	mc.totalCount++
	if data, ok := mc.m[key]; ok {
		// memory cache ignore this error
		node, _ := mc.lru.Access(data)
		if node == nil {
			return nil, nil
		}

		mc.hitCount++
		return node.Value.(CacheValue).Value, nil
	}
	return nil, nil
}

func (mc *memCache) GetHitInfo() (int64, int64) {
	mc.lock.RLock()
	defer mc.lock.RUnlock()

	return mc.hitCount, mc.totalCount
}

func newMemCache(maxSize int64, needCryptKey bool, ttl int64) Cache {
	if maxSize <= 0 {
		maxSize = DefaultMaxMemCacheSize
	}

	mc := &memCache{
		needCryptKey: needCryptKey,
		m:            make(map[interface{}]*lru.Node),
	}
	mc.lru = &lru.LRU{
		MaxSize:            maxSize,
		TTL:                ttl,
		AddNodeCallBack:    mc.addNodeCallback(),
		DeleteNodeCallBack: mc.deleteCallBack(),
	}
	return mc
}
