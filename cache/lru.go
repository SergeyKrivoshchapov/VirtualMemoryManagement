package cache

import "VirtualMemoryManagement/types/page"

type lruNode struct {
	page *page.Page
	prev *lruNode
	next *lruNode
}

type LRUCache struct {
	capacity int
	head     *lruNode
	tail     *lruNode
	pageMap  map[int]*lruNode
}

var _ Cache = (*LRUCache)(nil)

func NewLRU(capacity int) *LRUCache {
	lru := &LRUCache{
		capacity: capacity,
		pageMap:  make(map[int]*lruNode),
	}

	lru.head = &lruNode{}
	lru.tail = &lruNode{}
	lru.head.next = lru.tail
	lru.tail.prev = lru.head

	return lru
}

func (lru *LRUCache) Get(pageNumber int) *page.Page {
	if node, exists := lru.pageMap[pageNumber]; exists {
		lru.moveToFront(node)
		return node.page
	}
	return nil
}

func (lru *LRUCache) Put(p *page.Page) *page.Page {
	pageNumber := p.AbsoluteNumber

	if node, exists := lru.pageMap[pageNumber]; exists {
		node.page = p
		lru.moveToFront(node)
		return nil
	}

	var evicted *page.Page
	if len(lru.pageMap) >= lru.capacity {
		evicted = lru.evictLRU()
	}

	newNode := &lruNode{page: p}
	lru.pageMap[pageNumber] = newNode
	lru.addToFront(newNode)

	return evicted
}

// Contains checks if a page exists in the cache.
func (lru *LRUCache) Contains(pageNumber int) bool {
	_, exists := lru.pageMap[pageNumber]
	return exists
}

// Size returns the current number of pages in the cache.
func (lru *LRUCache) Size() int {
	return len(lru.pageMap)
}

// All returns a slice of all pages currently in the cache.
func (lru *LRUCache) All() []*page.Page {
	var result []*page.Page
	for _, node := range lru.pageMap {
		result = append(result, node.page)
	}
	return result
}

func (lru *LRUCache) moveToFront(node *lruNode) {
	lru.removeNode(node)
	lru.addToFront(node)
}

func (lru *LRUCache) addToFront(node *lruNode) {
	node.prev = lru.head
	node.next = lru.head.next
	lru.head.next.prev = node
	lru.head.next = node
}

func (lru *LRUCache) removeNode(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (lru *LRUCache) evictLRU() *page.Page {
	if lru.tail.prev == lru.head {
		return nil
	}

	node := lru.tail.prev
	lru.removeNode(node)
	delete(lru.pageMap, node.page.AbsoluteNumber)
	return node.page
}
