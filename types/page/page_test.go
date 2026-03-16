package page

import (
	"VirtualMemoryManagement/config"
	"testing"
	"time"
)

func TestPageNew(t *testing.T) {
	elemSize := 4
	p := New(42, elemSize)
	if p == nil {
		t.Fatal("Page should not be nil")
	}
	if p.AbsoluteNumber != 42 {
		t.Fatalf("Expected AbsoluteNumber 42, got %d", p.AbsoluteNumber)
	}
	if p.Dirty {
		t.Fatal("New page should not be dirty")
	}
	if p.WriteProtected {
		t.Fatal("New page should not be write protected")
	}
	if p.AccessCounter != 0 {
		t.Fatalf("Expected AccessCounter 0, got %d", p.AccessCounter)
	}
}
func TestPageDataSize(t *testing.T) {
	elemSize := 4
	p := New(0, elemSize)
	expectedSize := config.PageDataSize(elemSize)
	if len(p.Data()) != expectedSize {
		t.Fatalf("Expected data size %d, got %d", expectedSize, len(p.Data()))
	}
}
func TestPageBitmap(t *testing.T) {
	p := New(0, 4)
	bm := p.Bitmap()
	if bm == nil {
		t.Fatal("Bitmap should not be nil")
	}
}
func TestPageMakeDirty(t *testing.T) {
	p := New(0, 4)
	if p.Dirty {
		t.Fatal("Page should not be dirty initially")
	}
	p.MakeDirty()
	if !p.Dirty {
		t.Fatal("Page should be dirty after MakeDirty")
	}
}
func TestPageMarkAccessed(t *testing.T) {
	p := New(0, 4)
	oldTime := p.AccessTime
	time.Sleep(10 * time.Millisecond)
	p.MarkAccessed()
	if !p.AccessTime.After(oldTime) {
		t.Fatal("AccessTime should be updated")
	}
}
func TestPageMarkAccessedIncrementsCounter(t *testing.T) {
	p := New(0, 4)
	if p.AccessCounter != 0 {
		t.Fatalf("Expected AccessCounter 0, got %d", p.AccessCounter)
	}
	p.MarkAccessed()
	if p.AccessCounter != 1 {
		t.Fatalf("Expected AccessCounter 1, got %d", p.AccessCounter)
	}
	p.MarkAccessed()
	p.MarkAccessed()
	if p.AccessCounter != 3 {
		t.Fatalf("Expected AccessCounter 3, got %d", p.AccessCounter)
	}
}
func TestPageSetBit(t *testing.T) {
	p := New(0, 4)
	err := p.SetBit(5)
	if err != nil {
		t.Fatalf("Failed to set bit: %v", err)
	}
	isSet, _ := p.IsBitSet(5)
	if !isSet {
		t.Fatal("Bit 5 should be set")
	}
}
func TestPageSetBitMakesPageDirty(t *testing.T) {
	p := New(0, 4)
	if p.Dirty {
		t.Fatal("Page should not be dirty initially")
	}
	p.SetBit(0)
	if !p.Dirty {
		t.Fatal("Page should be dirty after SetBit")
	}
}
func TestPageNewWithData(t *testing.T) {
	elemSize := 4
	pageSize := config.PageDataSize(elemSize)
	bitmapData := make([]byte, config.BytesPerBitmap)
	bitmapData[0] = 0xFF
	pageData := make([]byte, pageSize)
	for i := range pageData {
		pageData[i] = byte(i % 256)
	}
	p, err := NewWithData(10, elemSize, bitmapData, pageData)
	if err != nil {
		t.Fatalf("Failed to create page with data: %v", err)
	}
	if p.AbsoluteNumber != 10 {
		t.Fatalf("Expected AbsoluteNumber 10, got %d", p.AbsoluteNumber)
	}
	if len(p.Data()) != pageSize {
		t.Fatalf("Expected data size %d, got %d", pageSize, len(p.Data()))
	}
	for i := 0; i < pageSize; i++ {
		if p.Data()[i] != pageData[i] {
			t.Fatalf("Data mismatch at position %d", i)
		}
	}
}
func TestPageNewWithDataInvalidBitmap(t *testing.T) {
	elemSize := 4
	pageSize := config.PageDataSize(elemSize)
	bitmapData := make([]byte, 5)
	pageData := make([]byte, pageSize)
	_, err := NewWithData(10, elemSize, bitmapData, pageData)
	if err == nil {
		t.Fatal("Expected error for invalid bitmap data")
	}
}
func TestPageWriteProtection(t *testing.T) {
	p := New(0, 4)
	if p.WriteProtected {
		t.Fatal("Page should not be write protected initially")
	}
	p.WriteProtected = true
	if !p.WriteProtected {
		t.Fatal("Page should be write protected")
	}
	p.WriteProtected = false
	if p.WriteProtected {
		t.Fatal("Page should not be write protected")
	}
}
func TestPageMultipleOperations(t *testing.T) {
	p := New(5, 4)

	p.MarkAccessed()
	p.MakeDirty()
	p.WriteProtected = true
	p.MarkAccessed()

	if p.AccessCounter != 3 {
		t.Fatalf("Expected AccessCounter 3, got %d", p.AccessCounter)
	}
	if !p.Dirty {
		t.Fatal("Page should be dirty")
	}
	if !p.WriteProtected {
		t.Fatal("Page should be write protected")
	}
	if p.AbsoluteNumber != 5 {
		t.Fatalf("Expected AbsoluteNumber 5, got %d", p.AbsoluteNumber)
	}
}

func TestPageDifferentElementSizes(t *testing.T) {
	testCases := []int{1, 4, 8, 16, 32}
	for _, elemSize := range testCases {
		p := New(0, elemSize)
		expectedSize := config.PageDataSize(elemSize)
		if len(p.Data()) != expectedSize {
			t.Fatalf("For elemSize %d: expected %d, got %d", elemSize, expectedSize, len(p.Data()))
		}
	}
}
func TestPageDataIsMutable(t *testing.T) {
	p := New(0, 4)
	data := p.Data()
	data[0] = 42
	data[1] = 100
	if p.Data()[0] != 42 || p.Data()[1] != 100 {
		t.Fatal("Page data should be mutable through Data() method")
	}
}
func TestPageDataSizeMethod(t *testing.T) {
	p := New(0, 4)
	expected := config.PageDataSize(4)
	actual := p.DataSize()
	if actual != expected {
		t.Fatalf("Expected DataSize %d, got %d", expected, actual)
	}
}
