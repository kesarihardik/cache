package cache

import (
	"fmt"
	"sync"
)

type lruNode[K comparable, V any] struct {
	key   K
	value V
	next  *lruNode[K, V]
	prev  *lruNode[K, V]
}

type lruCache[K comparable, V any] struct {
	capacity int
	nodeMap  map[K]*lruNode[K, V]
	head     *lruNode[K, V]
	tail     *lruNode[K, V]
	mu       sync.RWMutex
}

func NewLRUCache[K comparable, V any](capacity int) (Cache[K, V], error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("invalid cache size: %d", capacity)
	}
	head := &lruNode[K, V]{}
	tail := &lruNode[K, V]{}
	head.next = tail
	tail.prev = head

	return &lruCache[K, V]{
		capacity: capacity,
		nodeMap:  make(map[K]*lruNode[K, V]),
		head:     head,
		tail:     tail,
	}, nil
}

func (cache *lruCache[K, V]) addNode(node *lruNode[K, V]) {
	node.next = cache.head.next
	cache.head.next.prev = node
	node.prev = cache.head
	cache.head.next = node

	cache.nodeMap[node.key] = node
}

func (cache *lruCache[K, V]) deleteNode(node *lruNode[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev

	delete(cache.nodeMap, node.key)
}

func (cache *lruCache[K, V]) moveToFront(node *lruNode[K, V]) {
	cache.deleteNode(node)
	cache.addNode(node)
}

func (cache *lruCache[K, V]) Get(key K) (value V, ok bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	existingNode, exists := cache.nodeMap[key]
	if !exists {
		var zero V
		return zero, false
	}

	cache.moveToFront(existingNode)

	return existingNode.value, true
}

func (cache *lruCache[K, V]) Put(key K, value V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	existingNode, exists := cache.nodeMap[key]

	if exists {
		existingNode.value = value
		cache.moveToFront(existingNode)
	} else {
		if len(cache.nodeMap) == cache.capacity {
			cache.deleteNode(cache.tail.prev)
		}
		newNode := &lruNode[K, V]{key: key, value: value}
		cache.addNode(newNode)
	}
}

func (cache *lruCache[K, V]) Clear() {
	cache.nodeMap = map[K]*lruNode[K, V]{}
	cache.head = &lruNode[K, V]{}
	cache.tail = &lruNode[K, V]{}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head
}
