package fcache

import (
	"sync"
	"time"
)

// memData memory cache data
type memData struct {
	// key key of memory cache
	key string

	// value value of memory cache
	value []byte

	// accessTime timestamp of access time(second)
	accessTime int64

	// accessCount access count
	accessCount int64

	// expireTime expire time of memory cache data
	//expireTime int64

	// double linked list
	previous *memData
	next     *memData
}

// memCache memory cache.
type memCache struct {
	// m map of memory cache. the key is file full path name of request,
	// and map value is file data info.
	m map[interface{}]*memData

	// needCryptKey whether or not crypt key when Set and Get cache, default false.
	needCryptKey bool

	// double linked list header
	header *memData
	// double linked list tail
	tail *memData

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

// moveToHeader move cache node to header
func (mc *memCache) moveToHeader(md *memData) {
	if md != mc.header {
		if md == mc.tail {
			mc.tail = md.previous
			md.previous.next = nil
		} else {
			md.next.previous = md.previous
			md.previous.next = md.next
		}
		md.previous = nil
		md.next = mc.header
		mc.header.previous = md
		mc.header = md
	}
}

func (mc *memCache) newHeader(md *memData) {
	if mc.header != nil {
		md.next = mc.header
		mc.header.previous = md
		mc.header = md
	} else {
		mc.header = md
		mc.tail = md
	}
}

// eliminate eliminate cache of one of ten mc.maxSize capacity
func (mc *memCache) eliminate() {
	length := mc.maxSize / 10
	for mc.tail != nil && length > 0 {
		temp := mc.tail
		length -= int64(len(temp.value))
		mc.curSize -= int64(len(temp.value))

		mc.tail = temp.previous
		temp.previous = nil

		if mc.tail != nil {
			mc.tail.next = nil
		} else {
			mc.tail = nil
			mc.header = nil
		}
		delete(mc.m, temp.key)
	}
}

// Set set memory cache with key-value pair, and covered if key already exist.
func (mc *memCache) Set(key string, value []byte) {
	if mc.needCryptKey {
		key = MD5(key)
	}

	mc.lock.Lock()
	defer mc.lock.Unlock()

	if data, ok := mc.m[key]; ok {
		mc.moveToHeader(data)
		netCap := int64(len(value) - len(data.value))
		if mc.curSize+netCap > mc.maxSize {
			mc.eliminate()
		}
		mc.curSize += netCap
		data.value = value
		data.accessCount++
		data.accessTime = time.Now().Unix()
	} else {
		if mc.curSize+int64(len(value)) > mc.maxSize {
			mc.eliminate()
		}
		mc.curSize += int64(len(value))
		newData := &memData{
			key:         key,
			value:       value,
			accessTime:  time.Now().Unix(),
			accessCount: 1,
		}
		mc.m[key] = newData
		mc.newHeader(newData)
	}
}

// Get get memory cache with key, a error will be return if key is not exist.
func (mc *memCache) Get(key string) []byte {
	if mc.needCryptKey {
		key = MD5(key)
	}

	mc.lock.RLock()
	defer mc.lock.RUnlock()

	mc.totalCount++

	if data, ok := mc.m[key]; ok {
		mc.hitCount++

		// update access time and access count
		data.accessTime = time.Now().Unix()
		data.accessCount++

		mc.moveToHeader(data)

		return data.value
	}
	return nil
}

func newMemCache(maxSize int64, needCryptKey bool) Cache {
	if maxSize <= 0 {
		maxSize = DefaultMaxMemCacheSize
	}

	return &memCache{
		maxSize:      maxSize,
		needCryptKey: needCryptKey,
		m:            make(map[interface{}]*memData),
	}
}
