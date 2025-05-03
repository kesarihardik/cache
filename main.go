package main

import (
	"fmt"

	lfu "github.com/kesarihardik/cache/cacheUtil/lfu"
)

func main() {
	c, err := lfu.NewLFUCache[int, string](2)

	if err != nil {
		fmt.Print("not allocated")
	}

	fmt.Println(c.Get(2))

	c.Put(1, "first")
	c.Put(2, "second")
	fmt.Println(c.Get(2))
	c.Put(3, "third")
	fmt.Println(c.Get(1))
}
