package api

import (
	"VirtualMemoryManagement/tests/testutils"
	"strconv"
	"testing"
)

func TestVMCreateInt(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_int")
	result := VMCreate(filename, 1000, "int", 0)
	if result.Success != 1 {
		t.Fatalf("Expected success, got %d. Error: %s", result.Success, result.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got %d. Error: %s", openResult.Success, openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())
	if handle <= 0 {
		t.Fatal("Handle should be positive")
	}

	VMClose(handle)
}

func TestVMCreateChar(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_char")
	result := VMCreate(filename, 500, "char", 20)
	if result.Success != 1 {
		t.Fatalf("Expected success, got error: %s", result.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got error: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())
	VMClose(handle)
}

func TestVMCreateVarchar(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_varchar")
	result := VMCreate(filename, 2000, "varchar", 0)
	if result.Success != 1 {
		t.Fatalf("Expected success, got error: %s", result.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got error: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())
	VMClose(handle)
}

func TestVMCreateInvalidType(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_invalid")
	result := VMCreate(filename, 100, "invalid", 0)

	if result.Success != 0 {
		t.Fatal("Expected error for invalid type")
	}
}

func TestVMOpenNonexistent(t *testing.T) {
	result := VMOpen("/nonexistent/path")
	if result.Success != 0 {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestVMOpenTwiceAndReopenAfterClose(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_open_seq")
	createResult := VMCreate(filename, 100, "int", 0)
	if createResult.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", createResult.String())
	}

	first := VMOpen(filename)
	if first.Success != 1 {
		t.Fatalf("Expected first open success, got error: %s", first.String())
	}

	second := VMOpen(filename)
	if second.Success != 0 {
		t.Fatal("Expected error on second open of same file")
	}

	handle, _ := strconv.Atoi(first.String())
	VMClose(handle)

	third := VMOpen(filename)
	if third.Success != 1 {
		t.Fatalf("Expected open success after close, got error: %s", third.String())
	}

	newHandle, _ := strconv.Atoi(third.String())
	VMClose(newHandle)
}

func TestVMClose(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_close")
	createResult := VMCreate(filename, 100, "int", 0)
	if createResult.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got error: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())
	closeResult := VMClose(handle)
	if closeResult.Success != 1 {
		t.Fatalf("Expected success, got error: %s", closeResult.String())
	}
}

func TestVMCloseInvalidHandle(t *testing.T) {
	result := VMClose(9999)
	if result.Success != 0 {
		t.Fatal("Expected error for invalid handle")
	}
}

func TestVMRead(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_read")
	createResult := VMCreate(filename, 1000, "int", 0)
	if createResult.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got error: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())

	VMWrite(handle, 0, "42")

	result := VMRead(handle, 0)
	if result.Success != 1 {
		t.Fatalf("Expected success, got error: %s", result.String())
	}

	value := result.String()
	if value != "42" {
		t.Fatalf("Expected '42', got '%s'", value)
	}

	VMClose(handle)
}

func TestVMReadInvalidHandle(t *testing.T) {
	result := VMRead(9999, 0)
	if result.Success != 0 {
		t.Fatal("Expected error for invalid handle")
	}
}

func TestVMReadIndexOutOfRange(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_range")
	createResult := VMCreate(filename, 100, "int", 0)
	if createResult.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got error: %s", openResult.String())
	}

	handle, err := strconv.Atoi(openResult.String())

	if err != nil || handle <= 0 {
		t.Fatalf("Expected handle to be > 0, got %d", handle)
	}
	result := VMRead(handle, 200)
	if result.Success != 0 {
		t.Fatal("Expected error for index out of range")
	}

	VMClose(handle)
}

func TestVMWrite(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_write")
	createResult := VMCreate(filename, 1000, "int", 0)
	if createResult.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got error: %s", openResult.String())
	}
	handle, err := strconv.Atoi(openResult.String())
	if err != nil || handle <= 0 {
		t.Fatalf("Inalid handle: %v (%s)", err, openResult.String())
	}

	result := VMWrite(handle, 0, "100")
	if result.Success != 1 {
		t.Fatalf("Expected success, got error: %s", result.String())
	}

	VMClose(handle)
}

func TestVMWriteInvalidHandle(t *testing.T) {
	result := VMWrite(9999, 0, "42")
	if result.Success != 0 {
		t.Fatal("Expected error for invalid handle")
	}
}

func TestVMWriteIndexOutOfRange(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_write_range")
	createResult := VMCreate(filename, 100, "int", 0)
	if createResult.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Expected open success, got error: %s", openResult.String())
	}
	handle, _ := strconv.Atoi(openResult.String())

	result := VMWrite(handle, 200, "42")
	if result.Success != 0 {
		t.Fatal("Expected error for index out of range")
	}

	VMClose(handle)
}

func TestVMMultipleHandles(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename1 := testutils.TempFilePath(dir, "test_h1")
	filename2 := testutils.TempFilePath(dir, "test_h2")

	result1 := VMCreate(filename1, 500, "int", 0)
	if result1.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", result1.String())
	}
	result2 := VMCreate(filename2, 500, "int", 0)
	if result2.Success != 1 {
		t.Fatalf("Expected create success, got error: %s", result2.String())
	}

	open1 := VMOpen(filename1)
	open2 := VMOpen(filename2)

	handle1, _ := strconv.Atoi(open1.String())
	handle2, _ := strconv.Atoi(open2.String())

	if handle1 == handle2 {
		t.Fatal("Handles should be unique")
	}

	VMWrite(handle1, 0, "111")
	VMWrite(handle2, 0, "222")

	val1 := VMRead(handle1, 0)
	val2 := VMRead(handle2, 0)

	if val1.String() != "111" || val2.String() != "222" {
		t.Fatal("Handles should be independent")
	}

	VMClose(handle1)
	VMClose(handle2)
}

func TestVMCreateFileAlreadyExists(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_exists")
	res1 := VMCreate(filename, 100, "int", 0)
	if res1.Success != 1 {
		t.Fatalf("Expected first create success, got: %s", res1.String())
	}

	res2 := VMCreate(filename, 100, "int", 0)
	if res2.Success != 0 {
		t.Fatal("Expected error on second create for same file")
	}
}

func TestAPIWriteCloseReopenReadInt(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "api_persist_int")

	createRes := VMCreate(filename, 100, "int", 0)
	if createRes.Success != 1 {
		t.Fatalf("Create failed: %s", createRes.String())
	}

	open1 := VMOpen(filename)
	if open1.Success != 1 {
		t.Fatalf("First open failed: %s", open1.String())
	}
	handle1, err := strconv.Atoi(open1.String())
	if err != nil || handle1 <= 0 {
		t.Fatalf("Invalid handle: %v (%s)", err, open1.String())
	}

	writeRes := VMWrite(handle1, 5, "777")
	if writeRes.Success != 1 {
		t.Fatalf("Write failed: %s", writeRes.String())
	}

	closeRes := VMClose(handle1)
	if closeRes.Success != 1 {
		t.Fatalf("Close failed: %s", closeRes.String())
	}

	open2 := VMOpen(filename)
	if open2.Success != 1 {
		t.Fatalf("Second open failed: %s", open2.String())
	}
	handle2, err := strconv.Atoi(open2.String())
	if err != nil || handle2 <= 0 {
		t.Fatalf("Invalid handle2: %v (%s)", err, open2.String())
	}
	defer VMClose(handle2)

	readRes := VMRead(handle2, 5)
	if readRes.Success != 1 {
		t.Fatalf("Read after reopen failed: %s", readRes.String())
	}
	if readRes.String() != "777" {
		t.Fatalf("Expected 777, got %s", readRes.String())
	}
}
