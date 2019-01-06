package hpcache

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	MaxMemCacheSize    = 2 << 25 // 64M
	DiskCacheThreshold = 2 << 20 // 1M
	MaxDiskCacheSize   = 2 << 32 // 4G

	defaultDiskCachePath = "./cache/"
)

func MD5(src string) string {
	ctx := md5.New()
	ctx.Write([]byte(src))
	return hex.EncodeToString(ctx.Sum(nil))
}
