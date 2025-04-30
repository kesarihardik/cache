package main

import (
	"fmt"

	cache "github.com/kesarihardik/cache/cacheUtil"
)

func get[K comparable, V any](cache *cache.LRUCache[K, V], key K) {
	val, exists := cache.Get(key)
	if exists {
		fmt.Println(val)
	} else {
		fmt.Println("cache miss.")
	}
}

func main() {
	fmt.Print()
	cache, error := cache.NewLRUCache[int, string](2)

	if error != nil {
		fmt.Print(error)
	}

	get(cache, 2)

	cache.Put(1, "first")
	cache.Put(2, "second")

	get(cache, 1)

	cache.Put(3, "third")

	get(cache, 3)
}
