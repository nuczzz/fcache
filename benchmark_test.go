package fcache

import "testing"

func BenchmarkMemCache_Set(b *testing.B) {
	cache := NewMemCache(100, false)
	for i := 0; i < b.N; i++ {
		cache.Set("key", []byte("value"))
	}
}
