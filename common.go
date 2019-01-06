package fcache

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	//diskCacheThreshold = 2 << 20 // 1M

	DefaultMaxMemCacheSize  = 2 << 25 // 64M
	DefaultMaxDiskCacheSize = 2 << 32 // 4G
	DefaultDiskCacheDir     = "./cache/"
)

func MD5(src string) string {
	ctx := md5.New()
	ctx.Write([]byte(src))
	return hex.EncodeToString(ctx.Sum(nil))
}
