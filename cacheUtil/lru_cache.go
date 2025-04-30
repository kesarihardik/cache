package cacheUtil

import (
	"fmt"
	"sync"
)

type Node[K comparable, V any] struct {
	key   K
	value V
	next  *Node[K, V]
	prev  *Node[K, V]
}

type LRUCache[K comparable, V any] struct {
	capacity int
	store    map[K]*Node[K, V]
	head     *Node[K, V]
	tail     *Node[K, V]
	mu       sync.RWMutex
}

func NewLRUCache[K comparable, V any](capacity int) (*LRUCache[K, V], error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("can not initiate cache of size: %d", capacity)
	}
	head := &Node[K, V]{}
	tail := &Node[K, V]{}
	head.next = tail
	tail.prev = head

	return &LRUCache[K, V]{
		capacity: capacity,
		store:    make(map[K]*Node[K, V]),
		head:     head,
		tail:     tail,
	}, nil
}

func (cache *LRUCache[K, V]) addNodeToFront(node *Node[K, V]) {
	node.next = cache.head.next
	cache.head.next.prev = node
	node.prev = cache.head
	cache.head.next = node

	cache.store[node.key] = node
}

func (cache *LRUCache[K, V]) deleteNode(node *Node[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev

	delete(cache.store, node.key)
}

func (cache *LRUCache[K, V]) moveNodeToFront(node *Node[K, V]) {
	cache.deleteNode(node)
	cache.addNodeToFront(node)
}

func (cache *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	node, exists := cache.store[key]
	if !exists {
		var zero V
		return zero, false
	}

	cache.moveNodeToFront(node)

	return node.value, true
}

func (cache *LRUCache[K, V]) Put(key K, value V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	node, exists := cache.store[key]

	if exists {
		node.value = value
		cache.moveNodeToFront(node)
	} else {
		if len(cache.store) == cache.capacity {
			cache.deleteNode(cache.tail.prev)
		}
		newNode := &Node[K, V]{key: key, value: value}
		cache.addNodeToFront(newNode)
	}
}
