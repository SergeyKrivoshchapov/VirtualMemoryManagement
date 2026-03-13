package main

import (
	"VirtualMemoryManagement/storage"
	"VirtualMemoryManagement/types/array"
	"fmt"
	"os"
)

func main() {
	filename := "/tmp/test_debug.vm"
	os.Remove(filename)
	os.Remove(filename + ".varchar")

	fmt.Println("Creating PageFile...")
	pf := storage.NewPageFile(filename)

	fmt.Println("Calling Create...")
	err := pf.Create(100, array.TypeInt, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Success! File created.")
	pf.Close()

	info, _ := os.Stat(filename)
	if info != nil {
		fmt.Printf("File size: %d bytes\n", info.Size())
	}
}

