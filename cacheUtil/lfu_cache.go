package cacheUtil

import (
	"fmt"
	"sync"
)

type lfuNode[K comparable, V any] struct {
	key   K
	value V
	freq  int
	next  *lfuNode[K, V]
	prev  *lfuNode[K, V]
}

type nodePair[K comparable, V any] struct {
	head *lfuNode[K, V]
	tail *lfuNode[K, V]
}

type lfuCache[K comparable, V any] struct {
	capacity int
	freqMap  map[int]nodePair[K, V]
	nodeMap  map[K]*lfuNode[K, V]
	minFreq  int
	mu       sync.RWMutex
}

func NewLFUCache[K comparable, V any](cacheSize int) (Cache[K, V], error) {
	if cacheSize <= 0 {
		return nil, fmt.Errorf("can not initiate cache of size: %d", cacheSize)
	}

	return &lfuCache[K, V]{
		capacity: cacheSize,
		freqMap:  make(map[int]nodePair[K, V]),
		nodeMap:  make(map[K]*lfuNode[K, V]),
		minFreq:  0,
	}, nil
}

func (cache *lfuCache[K, V]) getHead(freq int) *lfuNode[K, V] {
	existingNodePair, exists := cache.freqMap[freq]
	if exists {
		return existingNodePair.head
	}

	newHead := &lfuNode[K, V]{}
	newTail := &lfuNode[K, V]{}
	newHead.next = newTail
	newTail.prev = newHead
	cache.freqMap[freq] = nodePair[K, V]{newHead, newTail}

	return cache.freqMap[freq].head
}

func (cache *lfuCache[K, V]) add(n *lfuNode[K, V]) {
	headNode := cache.getHead(n.freq)

	n.next = headNode.next
	n.prev = headNode
	headNode.next.prev = n
	headNode.next = n

	cache.nodeMap[n.key] = n
	cache.minFreq = n.freq
}

func (cache *lfuCache[K, V]) incrementFrequency(n *lfuNode[K, V]) {
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
		newNode := &lfuNode[K, V]{key: key, value: value, freq: 1}
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
	cache.nodeMap = map[K]*lfuNode[K, V]{}
	cache.minFreq = 0
}
