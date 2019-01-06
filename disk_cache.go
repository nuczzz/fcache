package hpcache

import (
	"github.com/goharbor/harbor/src/common/utils/log"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// diskData metadata of disk cache, but no value field of cache.
type diskData struct {
	// key disk cache key
	key string

	// size size of cache data
	size int

	accessTime int64

	accessCount int

	expireTime int64

	// double linked list
	previous *diskData
	next     *diskData
}

// diskCache disk cache
type diskCache struct {
	// dir directory of disk cache
	dir string

	// NeedCryptKey whether or not crypt key when Set and Get cache, default false.
	NeedCryptKey bool

	// m map of disk cache data, key is file name
	m map[string]*diskData

	// lock lock of disk cache
	lock sync.RWMutex

	// maxSize max size of memory cache data(byte).
	maxSize int

	// curSize current size of memory cache data.
	curSize int

	// hitCount hit cache count
	hitCount int

	// totalCount total count, contains hit count and missing count
	totalCount int

	header *diskData
	tail   *diskData
}

func (dc *diskCache) fileName(key string) string {
	return dc.dir + key
}

func (dc *diskCache) createFile(key string, value []byte) error {
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

func (dc *diskCache) eliminate() {
	length := dc.maxSize / 10
	for dc.tail != nil && length > 0 {
		temp := dc.tail
		length -= temp.size
		dc.curSize -= temp.size

		dc.tail = temp.previous
		temp.previous = nil

		if dc.tail != nil {
			dc.tail.next = nil
		} else {
			dc.tail = nil
			dc.header = nil
		}
		delete(dc.m, temp.key)

		if err := os.Remove(dc.fileName(temp.key)); err != nil && !os.IsNotExist(err) {
			log.Error(err)
		}
	}
}

func (dc *diskCache) Set(key string, value []byte) {
	if dc.NeedCryptKey {
		key = MD5(key)
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	// create file
	if err := dc.createFile(key, value); err != nil {
		log.Error(err)
		return
	}

	// change metadata
	if data, ok := dc.m[key]; ok {
		dc.moveToHeader(data)
		netCap := len(value) - data.size
		if dc.curSize+netCap > dc.maxSize {
			dc.eliminate()
		}
		dc.curSize += netCap

		data.accessCount++
		data.accessTime = time.Now().Unix()
	} else {
		if dc.curSize+len(value) > dc.maxSize {
			dc.eliminate()
		}
		dc.curSize += len(value)
		newData := &diskData{
			key:         key,
			size:        len(value),
			accessTime:  time.Now().Unix(),
			accessCount: 1,
		}
		dc.m[key] = newData
		dc.newHeader(newData)
	}
}

func (dc *diskCache) newHeader(dd *diskData) {
	if dc.header != nil {
		dd.next = dc.header
		dc.header.previous = dd
		dc.header = dd
	} else {
		dc.header = dd
		dc.tail = dd
	}
}

func (dc *diskCache) moveToHeader(dd *diskData) {
	if dd != dc.header {
		if dd == dc.tail {
			dc.tail = dd.previous
			dd.previous.next = nil
		} else {
			dd.next.previous = dd.previous
			dd.previous.next = dd.next
		}
		dd.previous = nil
		dd.next = dc.header
		dc.header.previous = dd
		dc.header = dd
	}
}

func (dc *diskCache) Get(key string) []byte {
	if dc.NeedCryptKey {
		key = MD5(key)
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	dc.totalCount++
	if data, ok := dc.m[key]; ok {
		dc.hitCount++

		data.accessTime = time.Now().Unix()
		data.accessCount++

		dc.moveToHeader(data)
		value, err := ioutil.ReadFile(dc.fileName(key))
		if err != nil {
			return nil
		}
		return value
	}

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
	dc := &diskCache{}
	if dc.dir != "" {
		if dc.dir[len(dc.dir)-1] != '/' {
			dc.dir += "/"
		}
	}
	return dc
}
