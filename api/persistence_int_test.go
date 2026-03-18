package api

import (
	"os"
	"strings"
	"testing"
)

// TestPersistenceWithInt - test the exact scenario from the user
func TestPersistenceWithInt(t *testing.T) {
	filename := "test_persistence_int.vm"

	// Cleanup
	os.Remove(filename)
	os.Remove(filename + ".varchar")
	defer os.Remove(filename)
	defer os.Remove(filename + ".varchar")

	// Step 1: Create array
	t.Log("Step 1: Create test(int)")
	result := VMCreate(filename, 100, "int", 0)
	if result.Success != 1 {
		t.Fatalf("Failed to create: %s", result.String())
	}

	// Step 2: Open file
	t.Log("Step 2: Open test")
	result = VMOpen(filename)
	if result.Success != 1 {
		t.Fatalf("Failed to open: %s", result.String())
	}
	handle1 := 1 // First handle

	// Step 3: Write 200 to index 1
	t.Log("Step 3: Input (1, 200)")
	result = VMWrite(handle1, 1, "200")
	if result.Success != 1 {
		t.Fatalf("Failed to write: %s", result.String())
	}

	// Step 4: Read back immediately
	t.Log("Step 4: Print (1) - immediate read")
	result = VMRead(handle1, 1)
	if result.Success != 1 {
		t.Fatalf("Failed to read: %s", result.String())
	}
	value := strings.TrimRight(string(result.Data[:]), "\x00")
	t.Logf("Immediate read result: %s", value)
	if value != "200" {
		t.Errorf("Expected 200, got %s", value)
	}

	// Step 5: Close file
	t.Log("Step 5: Exit (Close file)")
	result = VMClose(handle1)
	if result.Success != 1 {
		t.Fatalf("Failed to close: %s", result.String())
	}

	// Step 6: Simulate program restart - reopen file
	t.Log("Step 6: PROGRAM RESTART - Open test again")
	result = VMOpen(filename)
	if result.Success != 1 {
		t.Fatalf("Failed to reopen: %s", result.String())
	}
	handle2 := 2 // New handle after restart

	// Step 7: Read value after restart
	t.Log("Step 7: Print (1) - read after restart")
	result = VMRead(handle2, 1)
	if result.Success != 1 {
		t.Fatalf("Failed to read after restart: %s", result.String())
	}
	value = strings.TrimRight(string(result.Data[:]), "\x00")
	t.Logf("Read after restart result: %s", value)

	if value != "200" {
		t.Errorf("❌ PERSISTENCE FAILED: Expected 200 after restart, got %s", value)
	} else {
		t.Logf("✅ PERSISTENCE SUCCESS: Value 200 was preserved!")
	}

	VMClose(handle2)
}
