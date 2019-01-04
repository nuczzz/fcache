package hpcache

const (
	MaxMemCacheSize    = 2 << 25 // 64M
	DiskCacheThreshold = 2 << 20 // 1M
	MaxDiskCacheSize   = 2 << 32 // 4G

	DefaultDiskCachePath = "./cache"
)
