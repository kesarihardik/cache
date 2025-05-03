package lfu

import (
	"fmt"
	"sync"
)

type node[K comparable, V any] struct {
	key   K
	value V
	freq  int
	next  *node[K, V]
	prev  *node[K, V]
}

type nodePair[K comparable, V any] struct {
	head *node[K, V]
	tail *node[K, V]
}

type lfuCache[K comparable, V any] struct {
	capacity int
	freqMap  map[int]nodePair[K, V]
	nodeMap  map[K]*node[K, V]
	minFreq  int
	mu       sync.RWMutex
}

func NewLFUCache[K comparable, V any](cacheSize int) (*lfuCache[K, V], error) {
	if cacheSize <= 0 {
		return nil, fmt.Errorf("can not initiate cache of size: %d", cacheSize)
	}

	return &lfuCache[K, V]{
		capacity: cacheSize,
		freqMap:  make(map[int]nodePair[K, V]),
		nodeMap:  make(map[K]*node[K, V]),
		minFreq:  0,
	}, nil
}

func (cache *lfuCache[K, V]) getHead(freq int) *node[K, V] {
	existingNodePair, exists := cache.freqMap[freq]
	if exists {
		return existingNodePair.head
	}

	newHead := &node[K, V]{}
	newTail := &node[K, V]{}
	newHead.next = newTail
	newTail.prev = newHead
	cache.freqMap[freq] = nodePair[K, V]{newHead, newTail}

	return cache.freqMap[freq].head
}

func (cache *lfuCache[K, V]) add(n *node[K, V]) {
	headNode := cache.getHead(n.freq)

	n.next = headNode.next
	n.prev = headNode
	headNode.next.prev = n
	headNode.next = n

	cache.nodeMap[n.key] = n
	cache.minFreq = n.freq
}

func (cache *lfuCache[K, V]) incrementFrequency(n *node[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev

	delete(cache.nodeMap, n.key)

	head := cache.freqMap[n.freq].head
	if head.next == cache.freqMap[n.freq].tail {
		if cache.minFreq == n.freq {
			cache.minFreq++
		}
		delete(cache.freqMap, n.freq)
	}

	n.freq++

	headNode := cache.getHead(n.freq)
	n.next = headNode.next
	n.prev = headNode
	headNode.next.prev = n
	headNode.next = n

	cache.nodeMap[n.key] = n
}

func (cache *lfuCache[K, V]) Put(key K, value V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	existingNode, exists := cache.nodeMap[key]

	if exists {
		existingNode.value = value
		cache.incrementFrequency(existingNode)
	} else {
		if len(cache.nodeMap) == cache.capacity {
			lfuNode := cache.freqMap[cache.minFreq].tail.prev

			lfuNode.next.prev = lfuNode.prev
			lfuNode.prev.next = lfuNode.next

			delete(cache.nodeMap, lfuNode.key)
		}
		newNode := &node[K, V]{key: key, value: value, freq: 1}
		cache.add(newNode)
	}
}

func (cache *lfuCache[K, V]) Get(key K) (V, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	existingNode, exists := cache.nodeMap[key]

	if !exists {
		var zero V
		return zero, false
	}

	cache.incrementFrequency(existingNode)

	return existingNode.value, true
}

func (cache *lfuCache[K, V]) Clear() {
	cache.freqMap = make(map[int]nodePair[K, V])
	cache.nodeMap = map[K]*node[K, V]{}
	cache.minFreq = 0
}
