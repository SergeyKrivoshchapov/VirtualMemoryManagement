package cache

import (
	"VirtualMemoryManagement/types/page"
	"testing"
)

func TestLRUNewCache(t *testing.T) {
	lru := NewLRU(5)
	if lru == nil {
		t.Fatal("LRU cache should not be nil")
	}
	if lru.Size() != 0 {
		t.Fatalf("New cache should be empty, got size %d", lru.Size())
	}
}

func TestLRUPutAndGet(t *testing.T) {
	lru := NewLRU(3)

	p1 := page.New(1, 4)
	evicted := lru.Put(p1)

	if evicted != nil {
		t.Fatal("Should not evict when capacity not reached")
	}

	retrieved := lru.Get(1)
	if retrieved == nil {
		t.Fatal("Should retrieve the page")
	}
	if retrieved.AbsoluteNumber != 1 {
		t.Fatalf("Expected page 1, got %d", retrieved.AbsoluteNumber)
	}
}

func TestLRUEviction(t *testing.T) {
	lru := NewLRU(2)

	p1 := page.New(1, 4)
	p2 := page.New(2, 4)
	p3 := page.New(3, 4)

	lru.Put(p1)
	lru.Put(p2)

	if lru.Size() != 2 {
		t.Fatalf("Expected size 2, got %d", lru.Size())
	}

	evicted := lru.Put(p3)

	if evicted == nil {
		t.Fatal("Should evict a page")
	}
	if evicted.AbsoluteNumber != 1 {
		t.Fatalf("Expected page 1 to be evicted, got %d", evicted.AbsoluteNumber)
	}
	if lru.Size() != 2 {
		t.Fatalf("Size should remain 2, got %d", lru.Size())
	}
}

func TestLRUUpdateOnGet(t *testing.T) {
	lru := NewLRU(2)

	p1 := page.New(1, 4)
	p2 := page.New(2, 4)
	p3 := page.New(3, 4)

	lru.Put(p1)
	lru.Put(p2)

	lru.Get(1)

	evicted := lru.Put(p3)

	if evicted == nil {
		t.Fatal("Should evict a page")
	}
	if evicted.AbsoluteNumber != 2 {
		t.Fatalf("Expected page 2 to be evicted (p1 was accessed), got %d", evicted.AbsoluteNumber)
	}
}

func TestLRUUpdateOnPut(t *testing.T) {
	lru := NewLRU(2)

	p1 := page.New(1, 4)
	p2 := page.New(2, 4)

	lru.Put(p1)
	lru.Put(p2)

	p1Updated := page.New(1, 4)
	lru.Put(p1Updated)

	if lru.Size() != 2 {
		t.Fatalf("Size should remain 2, got %d", lru.Size())
	}

	p3 := page.New(3, 4)
	evicted := lru.Put(p3)

	if evicted == nil {
		t.Fatal("Should evict a page")
	}
	if evicted.AbsoluteNumber != 2 {
		t.Fatalf("Expected page 2 to be evicted, got %d", evicted.AbsoluteNumber)
	}
}

func TestLRUContains(t *testing.T) {
	lru := NewLRU(3)

	p1 := page.New(1, 4)
	lru.Put(p1)

	if !lru.Contains(1) {
		t.Fatal("Cache should contain page 1")
	}

	if lru.Contains(2) {
		t.Fatal("Cache should not contain page 2")
	}
}

func TestLRUSize(t *testing.T) {
	lru := NewLRU(5)

	if lru.Size() != 0 {
		t.Fatalf("Initial size should be 0, got %d", lru.Size())
	}

	for i := 1; i <= 5; i++ {
		p := page.New(i, 4)
		lru.Put(p)
		if lru.Size() != i {
			t.Fatalf("Expected size %d, got %d", i, lru.Size())
		}
	}
}

func TestLRUGetNonexistent(t *testing.T) {
	lru := NewLRU(3)

	p := lru.Get(999)
	if p != nil {
		t.Fatal("Get should return nil for nonexistent page")
	}
}

func TestLRUAll(t *testing.T) {
	lru := NewLRU(3)

	p1 := page.New(1, 4)
	p2 := page.New(2, 4)
	p3 := page.New(3, 4)

	lru.Put(p1)
	lru.Put(p2)
	lru.Put(p3)

	all := lru.All()
	if len(all) != 3 {
		t.Fatalf("Expected 3 pages, got %d", len(all))
	}

	pageNumbers := make(map[int]bool)
	for _, p := range all {
		pageNumbers[p.AbsoluteNumber] = true
	}

	if !pageNumbers[1] || !pageNumbers[2] || !pageNumbers[3] {
		t.Fatal("All pages should be present")
	}
}

func TestLRUCapacityOne(t *testing.T) {
	lru := NewLRU(1)

	p1 := page.New(1, 4)
	p2 := page.New(2, 4)

	lru.Put(p1)
	evicted := lru.Put(p2)

	if evicted == nil || evicted.AbsoluteNumber != 1 {
		t.Fatal("Page 1 should be evicted")
	}
	if lru.Size() != 1 {
		t.Fatalf("Size should be 1, got %d", lru.Size())
	}
}

func TestLRUMultipleEvictions(t *testing.T) {
	lru := NewLRU(2)

	pages := make([]*page.Page, 10)
	for i := 0; i < 10; i++ {
		pages[i] = page.New(i+1, 4)
	}

	evictedCount := 0
	for _, p := range pages {
		evicted := lru.Put(p)
		if evicted != nil {
			evictedCount++
		}
	}

	if evictedCount != 8 {
		t.Fatalf("Expected 8 evictions, got %d", evictedCount)
	}
	if lru.Size() != 2 {
		t.Fatalf("Final size should be 2, got %d", lru.Size())
	}
}

func TestLRULRUOrder(t *testing.T) {
	lru := NewLRU(3)

	p1 := page.New(1, 4)
	p2 := page.New(2, 4)
	p3 := page.New(3, 4)
	p4 := page.New(4, 4)

	lru.Put(p1)
	lru.Put(p2)
	lru.Put(p3)

	lru.Get(1)
	lru.Get(2)

	evicted := lru.Put(p4)

	if evicted == nil || evicted.AbsoluteNumber != 3 {
		t.Fatalf("Expected page 3 to be evicted, got %v", evicted)
	}
}
