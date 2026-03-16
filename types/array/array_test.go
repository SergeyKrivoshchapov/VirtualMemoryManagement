package array

import (
	"VirtualMemoryManagement/config"
	"testing"
)

func TestArrayInfoNewInt(t *testing.T) {
	info := NewInfo(1000, TypeInt, 0)

	if info.Size != 1000 {
		t.Fatalf("Expected size 1000, got %d", info.Size)
	}
	if info.Type != TypeInt {
		t.Fatalf("Expected type TypeInt, got %v", info.Type)
	}
	if info.ElementSize != 4 {
		t.Fatalf("Expected ElementSize 4 for int, got %d", info.ElementSize)
	}
}

func TestArrayInfoNewChar(t *testing.T) {
	info := NewInfo(500, TypeChar, 10)

	if info.Size != 500 {
		t.Fatalf("Expected size 500, got %d", info.Size)
	}
	if info.Type != TypeChar {
		t.Fatalf("Expected type TypeChar, got %v", info.Type)
	}
	if info.StringLength != 10 {
		t.Fatalf("Expected StringLength 10, got %d", info.StringLength)
	}
	if info.ElementSize != 10 {
		t.Fatalf("Expected ElementSize 10 for char, got %d", info.ElementSize)
	}
}

func TestArrayInfoNewVarchar(t *testing.T) {
	info := NewInfo(2000, TypeVarchar, 0)

	if info.Size != 2000 {
		t.Fatalf("Expected size 2000, got %d", info.Size)
	}
	if info.Type != TypeVarchar {
		t.Fatalf("Expected type TypeVarchar, got %v", info.Type)
	}
	if info.ElementSize != 4 {
		t.Fatalf("Expected ElementSize 4 for varchar pointers, got %d", info.ElementSize)
	}
}

func TestArrayInfoPageCount(t *testing.T) {
	testCases := []struct {
		size     int
		typ      Type
		expected int
	}{
		{128, TypeInt, 1},
		{129, TypeInt, 2},
		{256, TypeInt, 2},
		{257, TypeInt, 3},
	}

	for _, tc := range testCases {
		info := NewInfo(tc.size, tc.typ, 0)
		if info.PageCount != tc.expected {
			t.Fatalf("For size %d: expected %d pages, got %d", tc.size, tc.expected, info.PageCount)
		}
	}
}

func TestArrayInfoTypeString(t *testing.T) {
	testCases := []struct {
		typ      Type
		expected string
	}{
		{TypeInt, "int"},
		{TypeChar, "char"},
		{TypeVarchar, "varchar"},
	}

	for _, tc := range testCases {
		if tc.typ.String() != tc.expected {
			t.Fatalf("Expected %s, got %s", tc.expected, tc.typ.String())
		}
	}
}

func TestArrayInfoDifferentCharLengths(t *testing.T) {
	testCases := []int{1, 5, 10, 50, 100}

	for _, length := range testCases {
		info := NewInfo(1000, TypeChar, length)
		if info.ElementSize != length {
			t.Fatalf("For string length %d: expected ElementSize %d, got %d", length, length, info.ElementSize)
		}
		if info.StringLength != length {
			t.Fatalf("StringLength not preserved")
		}
	}
}

func TestArrayInfoPageCountBoundary(t *testing.T) {
	info1 := NewInfo(config.BitsPerPage, TypeInt, 0)
	info2 := NewInfo(config.BitsPerPage+1, TypeInt, 0)

	if info1.PageCount != 1 {
		t.Fatalf("Expected 1 page for size %d, got %d", config.BitsPerPage, info1.PageCount)
	}
	if info2.PageCount != 2 {
		t.Fatalf("Expected 2 pages for size %d, got %d", config.BitsPerPage+1, info2.PageCount)
	}
}

