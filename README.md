# fcache
File cache with LRU algorithm, contains memory cache and disk cache.

## Example
```
package main

import (
	"fmt"
    	
    	"github.com/nuczzz/fcache"
)

func main() {
	memCache := fcache.NewMemCache(10, false)
	memCache.Set("key1", []byte("123456789"))
	memCache.Set("key2", []byte("0"))
	memCache.Set("key3", []byte("1"))
	fmt.Println(memCache.Get("key3"))
	fmt.Println(memCache.Get("key1"))

	diskCache := fcache.NewDiskCache(10, false, "./cache")
	diskCache.Set("key1", []byte("123456789"))
	diskCache.Set("key2", []byte("0"))
	diskCache.Set("key3", []byte("1"))
	fmt.Println(diskCache.Get("key3"))
	fmt.Println(diskCache.Get("key1"))
}
```