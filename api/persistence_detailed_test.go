package api

import (
	"os"
	"strings"
	"testing"
)

// TestDetailedPersistence - detailed step-by-step test
func TestDetailedPersistence(t *testing.T) {
	filename := "test_detailed_persist.vm"

	// Cleanup before
	os.Remove(filename)
	os.Remove(filename + ".varchar")
	defer os.Remove(filename)
	defer os.Remove(filename + ".varchar")

	t.Log("=== STEP 1: Create array ===")
	result := VMCreate(filename, 100, "int", 0)
	if result.Success != 1 {
		t.Fatalf("Failed to create: %s", result.String())
	}
	t.Logf("File created successfully")

	t.Log("\n=== STEP 2: Open file (handle 1) ===")
	result = VMOpen(filename)
	if result.Success != 1 {
		t.Fatalf("Failed to open: %s", result.String())
	}
	t.Logf("File opened, handle: %s", result.String())

	t.Log("\n=== STEP 3: Write 20 to index 0 ===")
	result = VMWrite(1, 0, "20")
	if result.Success != 1 {
		t.Fatalf("Failed to write: %s", result.String())
	}
	t.Logf("Write completed")

	t.Log("\n=== STEP 4: Read index 0 immediately ===")
	result = VMRead(1, 0)
	if result.Success != 1 {
		t.Fatalf("Failed to read: %s", result.String())
	}
	readValue := strings.TrimRight(string(result.Data[:]), "\x00")
	t.Logf("Immediate read: %s", readValue)
	if readValue != "20" {
		t.Errorf("Expected 20, got %s", readValue)
	}

	t.Log("\n=== STEP 5: Get stats before close ===")
	result = VMStats(1)
	t.Logf("Stats:\n%s", result.String())

	t.Log("\n=== STEP 6: Close file (handle 1) ===")
	result = VMClose(1)
	if result.Success != 1 {
		t.Fatalf("Failed to close: %s", result.String())
	}
	t.Logf("File closed successfully")

	t.Log("\n=== STEP 7: Check file size after close ===")
	info, err := os.Stat(filename)
	if err != nil {
		t.Logf("Error getting file info: %v", err)
	} else {
		t.Logf("File size after close: %d bytes", info.Size())
	}

	t.Log("\n=== STEP 8: Open file again (handle 2) ===")
	result = VMOpen(filename)
	if result.Success != 1 {
		t.Fatalf("Failed to reopen: %s", result.String())
	}
	t.Logf("File reopened, handle: %s", result.String())

	t.Log("\n=== STEP 9: Get stats after reopen ===")
	result = VMStats(2)
	t.Logf("Stats:\n%s", result.String())

	t.Log("\n=== STEP 10: Read index 0 after reopen ===")
	result = VMRead(2, 0)
	if result.Success != 1 {
		t.Fatalf("Failed to read after reopen: %s", result.String())
	}
	readValue = strings.TrimRight(string(result.Data[:]), "\x00")
	t.Logf("Read after reopen: %s", readValue)

	if readValue != "20" {
		t.Errorf("❌ PERSISTENCE FAILED: Expected 20 after reopen, got %s", readValue)
	} else {
		t.Logf("✅ PERSISTENCE SUCCESS: Value 20 was preserved")
	}

	VMClose(2)
}
