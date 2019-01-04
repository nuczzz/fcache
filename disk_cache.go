package hpcache

type diskCache struct {
}

func (dc *diskCache) Set(key string, value []byte) {

}

func (dc *diskCache) Get(key string) []byte {
	return nil
}

func newDiskCache() Cache {
	return &diskCache{}
}
