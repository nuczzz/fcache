package fcache

import (
	"lru"
	"sync"
)

type MemValue struct {
	Value []byte
}

func (mv MemValue) Len() int64 {
	return int64(len(mv.Value))
}

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

	// maxSize max size of memory cache data(byte).
	maxSize int64

	// curSize current size of memory cache data.
	curSize int64

	// hitCount hit cache count
	hitCount int64

	// totalCount total count, contains hit count and missing count
	totalCount int64
}

func (mc *memCache) deleteCallBack() func(key interface{}) error {
	return func(key interface{}) error {
		mc.curSize -= mc.m[key].Length
		delete(mc.m, key)
		return nil
	}
}

// Set set memory cache with key-value pair, and covered if key already exist.
func (mc *memCache) Set(key string, value []byte) error {
	if mc.needCryptKey {
		key = MD5(key)
	}

	mc.lock.Lock()
	defer mc.lock.Unlock()

	if data, ok := mc.m[key]; ok {
		if err := mc.lru.Delete(data); err != nil {
			return err
		}
	}
	v := MemValue{Value: value}
	// memory cache ignore this error
	newNode, _ := mc.lru.AddNewNode(key, v)
	mc.curSize += newNode.Length
	mc.m[key] = newNode
	return nil
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
		return node.Value.(MemValue).Value, nil
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
		maxSize:      maxSize,
		needCryptKey: needCryptKey,
		m:            make(map[interface{}]*lru.Node),
	}
	link := &lru.LRU{
		MaxSize:            maxSize,
		TTL:                ttl,
		DeleteNodeCallBack: mc.deleteCallBack(),
	}
	mc.lru = link
	return mc
}
