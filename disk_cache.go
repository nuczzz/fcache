package fcache

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// diskData metadata of disk cache, but no value field of cache.
type diskData struct {
	// key disk cache key
	key string

	// size size of cache data
	size int64

	accessTime int64

	accessCount int64

	expireTime int64

	// double linked list
	previous *diskData
	next     *diskData
}

// diskCache disk cache
type diskCache struct {
	// dir directory of disk cache
	dir string

	// needCryptKey whether or not crypt key when Set and Get cache, default false.
	needCryptKey bool

	// m map of disk cache data, key is file name
	m map[string]*diskData

	// lock lock of disk cache
	lock sync.RWMutex

	// maxSize max size of memory cache data(byte).
	maxSize int64

	// curSize current size of memory cache data.
	curSize int64

	// hitCount hit cache count
	hitCount int64

	// totalCount total count, contains hit count and missing count
	totalCount int64

	header *diskData
	tail   *diskData

	// ttl time to live(second)
	ttl int64
}

func (dc *diskCache) exchange(node1, node2 *diskData) {
	pre1 := node1.previous
	pre2 := node2.previous
	next1 := node1.next
	next2 := node2.next

	if pre1 != nil {
		pre1.next = node2
	}
	node2.previous = pre1

	if pre2 != nil {
		pre2.next = node1

	}
	node1.previous = pre2

	if next1 != nil {
		next1.previous = node2

	}
	node2.next = next1

	if next2 != nil {
		next2.previous = node1
	}
	node1.next = next2

	if dc.header == node1 {
		dc.header = node2
	} else if dc.header == node2 {
		dc.header = node1
	}
	if dc.tail == node1 {
		dc.tail = node2
	} else if dc.tail == node2 {
		dc.tail = node1
	}
}

// sort sort double linked list by access time DESC
func (dc *diskCache) sort() {
	for i := dc.header; i != nil && i.next != nil; i = i.next {
		t := i
		for j := i.next; j != nil; j = j.next {
			if t.accessTime < j.accessTime {
				t = j
			}
		}
		if i != t {
			dc.exchange(i, t)
		}
	}
}

func (dc *diskCache) fileName(key string) string {
	return dc.dir + key
}

func (dc *diskCache) createFile(key string, value []byte) error {
	dc.initDir()

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
			log.Fatal(err)
		}
	}
}

func (dc *diskCache) Set(key string, value []byte) {
	if dc.needCryptKey {
		key = MD5(key)
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	// create file
	if err := dc.createFile(key, value); err != nil {
		log.Fatal(err)
		return
	}

	now := time.Now().Unix()
	// change metadata
	if data, ok := dc.m[key]; ok {
		dc.moveToHeader(data)
		netCap := int64(len(value)) - data.size
		if dc.curSize+netCap > dc.maxSize {
			dc.eliminate()
		}
		dc.curSize += netCap

		data.accessCount++
		data.accessTime = now
		data.expireTime = now + dc.ttl
	} else {
		if dc.curSize+int64(len(value)) > dc.maxSize {
			dc.eliminate()
		}
		dc.curSize += int64(len(value))
		newData := &diskData{
			key:        key,
			size:       int64(len(value)),
			accessTime: now,
			expireTime: now + dc.ttl,
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
	if dc.needCryptKey {
		key = MD5(key)
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	dc.totalCount++
	if data, ok := dc.m[key]; ok {
		dc.hitCount++

		data.accessTime = time.Now().UnixNano()
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
func (dc *diskCache) init() error {
	if err := dc.initDir(); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(dc.dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		fileName := dc.fileName(file.Name())
		if file.IsDir() {
			continue
		}

		fi, err := GetFileTime(fileName)
		if err != nil {
			return err
		}
		data := &diskData{
			key:        file.Name(),
			size:       file.Size(),
			accessTime: fi.AccessTime / 1e6,
		}
		// double linked list init
		dc.newHeader(data)

		// disk cache metadata init
		dc.m[file.Name()] = data
		if dc.curSize+file.Size() > dc.maxSize {
			if err := os.Remove(fileName); err != nil {
				log.Println(err)
			}
		} else {
			dc.curSize += file.Size()
		}
	}
	dc.sort()

	return nil
}

func (dc *diskCache) GetHitInfo() (int64, int64) {
	dc.lock.RLock()
	defer dc.lock.RUnlock()

	return dc.hitCount, dc.totalCount
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
	/*return &diskCache{
		maxSize:      maxSize,
		needCryptKey: needCryptKey,
		dir:          cacheDir,
		m:            make(map[string]*diskData),
		ttl:          ttl,
	}*/
	return nil
}
