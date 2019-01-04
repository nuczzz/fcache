package hpcache

// Cache cache interface definition.The cache can be memory cache,
// disk cache or net cache. We implementation cache with LRU algorithm.
type Cache interface {
	// Set set cache with key-value pair.
	Set(key string, value []byte)

	// Get get cache with key, nil will be return if key is not exist.
	Get(key string) []byte
}

func NewMemCache() Cache {
	return newMemCache()
}

func NewDiskCache() Cache {
	return newDiskCache()
}
