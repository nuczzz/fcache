package hpcache

import (
	"io/ioutil"
	"sync"
)

type diskData struct {
	// key disk cache key
	key string

	// value value of disk cache
	value []byte

	accessTime int64

	expireTime int64

	// double linked list
	previous *diskData
	next     *diskData
}

// diskCache disk cache
type diskCache struct {
	// dir directory of disk cache
	dir string

	// m map of disk cache data, key is file name
	m map[string]*diskData

	// lock lock of disk cache
	lock sync.RWMutex

	// hitCount hit cache count
	hitCount int

	// totalCount total count, contains hit count and missing count
	totalCount int
}

func (dc *diskCache) fileName(key string) string {
	return dc.dir + key
}

func (dc *diskCache) Set(key string, value []byte) {
	key = MD5(key)

	dc.lock.Lock()
	defer dc.lock.Unlock()

}

func (dc *diskCache) Get(key string) []byte {
	key = MD5(key)

	dc.lock.Lock()
	defer dc.lock.Unlock()

	return nil
}

// init read disk cache info when create new disk cache
func (dc *diskCache) init() error {
	files, err := ioutil.ReadDir(dc.dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		//todo
	}

	return nil
}

func newDiskCache() Cache {
	return &diskCache{}
}
