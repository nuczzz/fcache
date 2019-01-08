package fcache

// Cache cache interface definition.The cache can be memory cache,
// disk cache or net cache. We implementation cache with LRU algorithm.
type Cache interface {
	// Set set cache with key-value pair.
	Set(key string, value []byte)

	// Get get cache with key, nil will be return if key is not exist.
	Get(key string) []byte

	// GetHitInfo get cache hit info, return the count of hit visitor and the count of total visitor
	GetHitInfo() (hitCount, totalCount int64)
}

func NewMemCache(maxSize int64, needCryptKey bool, ttl ...int64) Cache {
	if len(ttl) > 0 {
		return newMemCache(maxSize, needCryptKey, ttl[0])
	}
	return newMemCache(maxSize, needCryptKey, 0)
}

func NewDiskCache(maxSize int64, needCryptKey bool, cacheDir string, ttl ...int64) Cache {
	if len(ttl) > 0 {
		return newDiskCache(maxSize, needCryptKey, cacheDir, ttl[0])
	}
	return newDiskCache(maxSize, needCryptKey, cacheDir, 0)
}
