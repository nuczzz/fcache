package hpcache

type diskCache struct {

}

func (dc *diskCache) SetCache(key string, value []byte) error {
	return nil
}

func (dc *diskCache) GetCache(key string) (value []byte, err error) {
	return nil, nil
}


func newDiskCache() Cache {
	return &diskCache{}
}
