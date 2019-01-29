package fcache

import (
	"crypto/md5"
	"encoding/hex"
	"unsafe"
)

const (
	//diskCacheThreshold = 2 << 20 // 1M

	DefaultMaxMemCacheSize  = 2 << 25 // 64M
	DefaultMaxDiskCacheSize = 2 << 32 // 4G
	DefaultDiskCacheDir     = "./cache/"
)

func String2Bytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func MD5(src string) string {
	ctx := md5.New()
	ctx.Write(String2Bytes(src))
	return hex.EncodeToString(ctx.Sum(nil))
}
