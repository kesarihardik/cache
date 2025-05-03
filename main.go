package main

import (
	"fmt"

	cache "github.com/kesarihardik/cache/cacheUtil"
)

func main() {
	c, err := cache.NewLFUCache[int, string](2)

	if err != nil {
		fmt.Print("not allocated")
	}

	fmt.Println(c.Get(2))

	c.Put(1, "first")
	c.Put(2, "second")

	fmt.Println(c.Get(1))

	c.Put(3, "third")
	fmt.Println(c.Get(1))
	fmt.Println(c.Get(2))
	fmt.Println(c.Get(3))
}
