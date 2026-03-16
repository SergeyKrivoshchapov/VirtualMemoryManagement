package bitmap

import (
	"VirtualMemoryManagement/config"
	"testing"
)

func TestBitMapNew(t *testing.T) {
	bm := New()
	if bm == nil {
		t.Fatal("BitMap should not be nil")
	}
}
func TestBitMapSetAndIsSet(t *testing.T) {
	bm := New()
	for pos := 0; pos < config.BitsPerPage; pos++ {
		err := bm.Set(pos)
		if err != nil {
			t.Fatalf("Failed to set bit at position %d: %v", pos, err)
		}
		isSet, err := bm.IsSet(pos)
		if err != nil {
			t.Fatalf("Failed to check bit at position %d: %v", pos, err)
		}
		if !isSet {
			t.Fatalf("Bit at position %d should be set", pos)
		}
	}
}
func TestBitMapSetInvalidPositionNegative(t *testing.T) {
	bm := New()
	err := bm.Set(-1)
	if err == nil {
		t.Fatal("Expected error for negative position")
	}
}
func TestBitMapSetInvalidPositionTooLarge(t *testing.T) {
	bm := New()
	err := bm.Set(config.BitsPerPage)
	if err == nil {
		t.Fatal("Expected error for position >= BitsPerPage")
	}
}
func TestBitMapIsSetInvalidPositionNegative(t *testing.T) {
	bm := New()
	_, err := bm.IsSet(-1)
	if err == nil {
		t.Fatal("Expected error for negative position")
	}
}
func TestBitMapIsSetInvalidPositionTooLarge(t *testing.T) {
	bm := New()
	_, err := bm.IsSet(config.BitsPerPage)
	if err == nil {
		t.Fatal("Expected error for position >= BitsPerPage")
	}
}
func TestBitMapBytes(t *testing.T) {
	bm := New()
	bm.Set(0)
	bm.Set(7)
	bm.Set(15)
	data := bm.Bytes()
	if len(data) != config.BytesPerBitmap {
		t.Fatalf("Expected %d bytes, got %d", config.BytesPerBitmap, len(data))
	}
}
func TestBitMapFromBytes(t *testing.T) {
	bm1 := New()
	bm1.Set(5)
	bm1.Set(13)
	bm1.Set(100)
	data := bm1.Bytes()
	bm2 := New()
	err := bm2.FromBytes(data)
	if err != nil {
		t.Fatalf("Failed to load from bytes: %v", err)
	}
	for pos := 0; pos < config.BitsPerPage; pos++ {
		isSet1, _ := bm1.IsSet(pos)
		isSet2, _ := bm2.IsSet(pos)
		if isSet1 != isSet2 {
			t.Fatalf("Bit %d mismatch: %v vs %v", pos, isSet1, isSet2)
		}
	}
}
func TestBitMapFromBytesInvalidLength(t *testing.T) {
	bm := New()
	invalidData := make([]byte, 5)
	err := bm.FromBytes(invalidData)
	if err == nil {
		t.Fatal("Expected error for invalid data length")
	}
}
func TestBitMapMultipleOperations(t *testing.T) {
	bm := New()
	positions := []int{0, 10, 20, 50, 100, 127}
	for _, pos := range positions {
		bm.Set(pos)
	}
	for i := 0; i < config.BitsPerPage; i++ {
		isSet, _ := bm.IsSet(i)
		expectedSet := false
		for _, p := range positions {
			if p == i {
				expectedSet = true
				break
			}
		}
		if isSet != expectedSet {
			t.Fatalf("Position %d: expected %v, got %v", i, expectedSet, isSet)
		}
	}
}
func TestBitMapBoundaryBits(t *testing.T) {
	bm := New()
	testCases := []int{0, 7, 8, 15, config.BitsPerPage - 1}
	for _, pos := range testCases {
		bm.Set(pos)
		isSet, err := bm.IsSet(pos)
		if err != nil {
			t.Fatalf("Error at position %d: %v", pos, err)
		}
		if !isSet {
			t.Fatalf("Boundary bit at position %d should be set", pos)
		}
	}
}
