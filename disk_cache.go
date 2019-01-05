package hpcache

import (
	"io/ioutil"
)

type diskData struct {
}

type diskCache struct {
	// dir directory of disk cache
	dir string

	// m map of disk cache data
	m map[string]*diskData
}

func (dc *diskCache) Set(key string, value []byte) {

}

func (dc *diskCache) Get(key string) []byte {
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
