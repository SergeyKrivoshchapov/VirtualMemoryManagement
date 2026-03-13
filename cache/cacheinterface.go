// Package cache provides caching interfaces and implementations for page management.
package cache

import "VirtualMemoryManagement/types/page"

// Cache defines the interface for page caching with LRU eviction.
type Cache interface {
	Get(pageNumber int) *page.Page
	Put(p *page.Page) *page.Page
	Contains(pageNumber int) bool
	Size() int
	All() []*page.Page
}
