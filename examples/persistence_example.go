// Manual Testing Examples for Virtual Memory Management
// These are example programs to manually test the persistence functionality

package main

import (
	"fmt"
	"os"
	"strings"

	"VirtualMemoryManagement/api"
)

func examplePersistenceTest() {
	testFile := "persistence_test.vm"

	// Cleanup
	os.Remove(testFile)
	os.Remove(testFile + ".varchar")

	fmt.Println("Step 1: Create array")
	result := api.VMCreate(testFile, 100, "int", 0)
	fmt.Printf("  Result: Success=%d, %s\n", result.Success, result.String())

	fmt.Println("\nStep 2: Open file")
	result = api.VMOpen(testFile)
	fmt.Printf("  Result: Success=%d, Handle=%s\n", result.Success, result.String())

	fmt.Println("\nStep 3: Write 20 to index 0")
	result = api.VMWrite(1, 0, "20")
	fmt.Printf("  Result: Success=%d, %s\n", result.Success, result.String())

	fmt.Println("\nStep 4: Read immediately")
	result = api.VMRead(1, 0)
	value := strings.TrimRight(string(result.Data[:]), "\x00")
	fmt.Printf("  Result: Success=%d, Value=%s\n", result.Success, value)

	fmt.Println("\nStep 5: Close file")
	result = api.VMClose(1)
	fmt.Printf("  Result: Success=%d, %s\n", result.Success, result.String())

	fmt.Println("\nStep 6: Reopen file")
	result = api.VMOpen(testFile)
	fmt.Printf("  Result: Success=%d, Handle=%s\n", result.Success, result.String())

	fmt.Println("\nStep 7: Read after reopen")
	result = api.VMRead(2, 0)
	value = strings.TrimRight(string(result.Data[:]), "\x00")
	fmt.Printf("  Result: Success=%d, Value=%s\n", result.Success, value)

	if value == "20" {
		fmt.Println("\n✅ SUCCESS: Persistence works correctly!")
	} else {
		fmt.Printf("\n❌ FAILURE: Expected '20', got '%s'\n", value)
	}

	api.VMClose(2)

	// Cleanup
	os.Remove(testFile)
	os.Remove(testFile + ".varchar")
}
