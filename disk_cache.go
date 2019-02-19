package fcache

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/nuczzz/lru"
	"sync/atomic"
	"time"
)

// diskCache disk cache
type diskCache struct {
	// dir directory of disk cache
	dir string

	// needCryptKey whether or not crypt key when Set and Get cache, default false.
	needCryptKey bool

	// m map of disk cache data, key is file name
	m map[interface{}]*lru.Node

	// lock lock of disk cache
	lock sync.RWMutex

	// hitCount hit cache count
	hitCount int64

	// totalCount total count, contains hit count and missing count
	totalCount int64

	lru *lru.LRU
}

func (dc *diskCache) fileName(key string) string {
	return dc.dir + key
}

func (dc *diskCache) createFile(key string, value []byte) error {
	if err := dc.initDir(); err != nil {
		return err
	}

	fd, err := os.Create(dc.fileName(key))
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err = fd.Write(value); err != nil {
		return err
	}
	return nil
}

func (dc *diskCache) addNodeCallback() func(node *lru.Node) {
	return func(node *lru.Node) {
		dc.m[node.Key] = node
	}
}

func (dc *diskCache) Set(key string, value []byte, extra ...interface{}) error {
	if dc.needCryptKey {
		key = MD5(key)
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	v := CacheValue{Value: value}
	if data, ok := dc.m[key]; ok {
		if err := dc.lru.Replace(data, v); err != nil {
			return err
		}
	}
	return dc.lru.AddNewNode(key, v, extra...)
}

func (dc *diskCache) Get(key string) (value []byte, extra interface{}, err error) {
	if dc.needCryptKey {
		key = MD5(key)
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	atomic.AddInt64(&dc.totalCount, 1)
	if data, ok := dc.m[key]; ok {
		node, err := dc.lru.Access(data)
		if err != nil {
			return nil, nil, err
		}
		if node == nil {
			return nil, nil, nil
		}

		atomic.AddInt64(&dc.hitCount, 1)
		return node.Value.(CacheValue).Value, node.Extra, nil
	}

	return nil, nil, nil
}

// initDir check disk cache directory exist or not, create it if not exist.
func (dc *diskCache) initDir() error {
	fd, err := os.Open(dc.dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dc.dir, 0755)
	}
	fd.Close()
	return nil
}

// init read disk cache file info when create new disk cache
// todo: how to store and read EXTRA field??
func (dc *diskCache) init() error {
	files, err := ioutil.ReadDir(dc.dir)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, file := range files {
		fileName := dc.fileName(file.Name())
		if file.IsDir() {
			continue
		}

		fi, err := GetFileTime(fileName)
		if err != nil {
			return err
		}

		node := &lru.Node{
			Key:        file.Name(),
			Length:     file.Size(),
			AccessTime: fi.AccessTime / 1e9,
		}
		if dc.lru.TTL > 0 {
			if fi.AccessTime/1e9+dc.lru.TTL <= time.Now().Unix() {
				os.Remove(dc.dir + file.Name())
				continue
			} else {
				node.SetExpire(now + dc.lru.TTL)
			}
		}

		dc.lru.Add(node)
		dc.m[node.Key] = node
	}
	return nil
}

func (dc *diskCache) GetHitInfo() (int64, int64) {
	dc.lock.RLock()
	defer dc.lock.RUnlock()

	return atomic.LoadInt64(&dc.hitCount), atomic.LoadInt64(&dc.totalCount)
}

func (dc *diskCache) deleteCallBack() func(key interface{}) error {
	return func(key interface{}) error {
		if err := os.Remove(dc.fileName(key.(string))); err != nil {
			return err
		}
		delete(dc.m, key)
		return nil
	}
}

func (dc *diskCache) setValue() func(key, value interface{}) error {
	return func(key, value interface{}) error {
		return dc.createFile(key.(string), value.(CacheValue).Value)
	}
}

func (dc *diskCache) getValue() func(interface{}) (lru.Value, error) {
	return func(key interface{}) (lru.Value, error) {
		value, err := ioutil.ReadFile(dc.fileName(key.(string)))
		return CacheValue{Value: value}, err
	}
}

func (dc *diskCache) Clear(key string) error {
	if dc.needCryptKey {
		key = MD5(key)
	}

	dc.lock.RLock()
	defer dc.lock.RUnlock()

	if data, ok := dc.m[key]; ok {
		return dc.lru.Delete(data)
	}
	return nil
}

func (dc *diskCache) ClearAll() error {
	dc.lock.RLock()
	defer dc.lock.RUnlock()

	var err error
	for _, node := range dc.lru.Traversal() {
		if err = dc.lru.Delete(node); err != nil {
			return err
		}
	}
	return nil
}

func newDiskCache(maxSize int64, needCryptKey bool, cacheDir string, ttl int64) Cache {
	if maxSize <= 0 {
		maxSize = DefaultMaxDiskCacheSize
	}
	if cacheDir == "" {
		cacheDir = DefaultDiskCacheDir
	}
	if cacheDir[len(cacheDir)-1] != '/' {
		cacheDir += "/"
	}
	dc := &diskCache{
		needCryptKey: needCryptKey,
		dir:          cacheDir,
		m:            make(map[interface{}]*lru.Node),
	}
	dc.lru = &lru.LRU{
		MaxSize:            maxSize,
		TTL:                ttl,
		AddNodeCallBack:    dc.addNodeCallback(),
		DeleteNodeCallBack: dc.deleteCallBack(),
		SetValue:           dc.setValue(),
		GetValue:           dc.getValue(),
	}
	return dc
}
