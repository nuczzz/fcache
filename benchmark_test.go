package fcache

import (
	"testing"
)

// go test -run=^^$ -bench=^BenchmarkMemCacheSet$ -benchmem
func BenchmarkMemCacheSet(b *testing.B) {
	cache := NewMemCache(100, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", []byte("value"))
	}
}

// go test -run=^^$ -bench=^BenchmarkMemCacheGet$ -benchmem
func BenchmarkMemCacheGet(b *testing.B) {
	cache := NewMemCache(100, false)
	cache.Set("key", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

// go test -run=^^$ -bench=^BenchmarkDiskCacheSet$ -benchmem
func BenchmarkDiskCacheSet(b *testing.B) {
	cache := NewDiskCache(100, false, "./cache")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", []byte("value"))
	}
}

// go test -run=^^$ -bench=^BenchmarkDiskCacheGet$ -benchmem
func BenchmarkDiskCacheGet(b *testing.B) {
	cache := NewDiskCache(100, false, "./cache")
	cache.Set("key", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}
