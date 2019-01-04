package hpcache

type memCache struct {
}

func (mc *memCache) SetCache(key string, value []byte) error {
	return nil
}

func (mc *memCache) GetCache(key string) (value []byte, err error) {
	return nil, nil
}

func newMemCache() Cache {
	return &memCache{}
}
