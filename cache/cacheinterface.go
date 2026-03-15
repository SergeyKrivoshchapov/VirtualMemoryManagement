package cache

import "VirtualMemoryManagement/types/page"

type Cache interface {
	Get(pageNumber int) *page.Page
	Put(p *page.Page) *page.Page
	Contains(pageNumber int) bool
	Size() int
	All() []*page.Page
}
