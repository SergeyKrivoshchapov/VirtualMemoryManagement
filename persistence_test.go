package main

import (
	"VirtualMemoryManagement/api"
	"os"
	"strings"
	"testing"
)

func TestPersistenceIssue(t *testing.T) {
	filename := "test_persist.vm"

	defer os.Remove(filename)
	defer os.Remove(filename + ".varchar")

	result := api.VMCreate(filename, 100, "int", 0)
	if result.Success != 1 {
		t.Fatalf("Failed to create: %s", result.String())
	}

	result = api.VMOpen(filename)
	if result.Success != 1 {
		t.Fatalf("Failed to open: %s", result.String())
	}

	result = api.VMWrite(1, 0, "20")
	if result.Success != 1 {
		t.Fatalf("Failed to write: %s", result.String())
	}

	result = api.VMRead(1, 0)
	if result.Success != 1 {
		t.Fatalf("Failed to read immediately: %s", result.String())
	}
	if strings.TrimRight(string(result.Data[:]), "\x00") != "20" {
		t.Errorf("Immediate read failed: expected 20, got %s", result.String())
	}

	result = api.VMClose(1)
	if result.Success != 1 {
		t.Fatalf("Failed to close: %s", result.String())
	}

	if result.Success != 1 {
		t.Fatalf("Failed to reopen: %s", result.String())
	}

	result = api.VMRead(2, 0)
	if result.Success != 1 {
		t.Fatalf("Failed to read after reopen: %s", result.String())
	}
	if strings.TrimRight(string(result.Data[:]), "\x00") != "20" {
		t.Errorf("PERSISTENCE ISSUE DETECTED: expected 20 after reopen, got %s", result.String())
	}

	api.VMClose(2)
}

func TestMultipleWritesPersist(t *testing.T) {
	filename := "test_multi_persist.vm"

	defer os.Remove(filename)
	defer os.Remove(filename + ".varchar")

	result := api.VMCreate(filename, 100, "int", 0)
	if result.Success != 1 {
		t.Fatalf("Failed to create: %s", result.String())
	}

	result = api.VMOpen(filename)
	if result.Success != 1 {
		t.Fatalf("Failed to open: %s", result.String())
	}

	testCases := []struct {
		index int
		value string
		name  string
	}{
		{0, "10", "first"},
		{5, "50", "middle"},
		{99, "100", "last"},
	}

	for _, tc := range testCases {
		result = api.VMWrite(3, tc.index, tc.value)
		if result.Success != 1 {
			t.Fatalf("Failed to write %s: %s", tc.name, result.String())
		}
	}

	api.VMClose(3)

	result = api.VMOpen(filename)
	if result.Success != 1 {
		t.Fatalf("Failed to reopen: %s", result.String())
	}

	for _, tc := range testCases {
		result = api.VMRead(4, tc.index)
		if result.Success != 1 {
			t.Fatalf("Failed to read %s after reopen: %s", tc.name, result.String())
		}
		if strings.TrimRight(string(result.Data[:]), "\x00") != tc.value {
			t.Errorf("Value mismatch for %s: expected %s, got %s", tc.name, tc.value, result.String())
		}
	}

	api.VMClose(4)
}
