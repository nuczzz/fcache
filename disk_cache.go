package hpcache

type diskCache struct {
}

func (dc *diskCache) Set(key interface{}, value interface{}) {

}

func (dc *diskCache) Get(key interface{}) interface{} {
	return nil
}

func newDiskCache() Cache {
	return &diskCache{}
}
