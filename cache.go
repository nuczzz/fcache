package hpcache

type Cache interface {
	SetCache(key string, value []byte) error
	GetCache(key string) (value []byte, err error)
}

func NewMemCache() Cache {
	return newMemCache()
}

func NewDiskCache() Cache {
	return newDiskCache()
}
