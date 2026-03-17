package api

import (
	"VirtualMemoryManagement/tests/testutils"
	"strconv"
	"strings"
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

func TestAPICharArrayReadWritePersist(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "api_char_persist")
	stringLen := 10

	createRes := VMCreate(filename, 20, "char", stringLen)
	if createRes.Success != 1 {
		t.Fatalf("Create failed: %s", createRes.String())
	}

	open1 := VMOpen(filename)
	if open1.Success != 1 {
		t.Fatalf("Open failed: %s", open1.String())
	}
	handle1, _ := strconv.Atoi(open1.String())
	defer VMClose(handle1)

	writeRes1 := VMWrite(handle1, 0, "ab")
	if writeRes1.Success != 1 {
		t.Fatalf("Write short string failed: %s", writeRes1.String())
	}

	writeRes2 := VMWrite(handle1, 1, "123456789012345")
	if writeRes2.Success != 1 {
		t.Fatalf("Write long string failed: %s", writeRes2.String())
	}

	readRes1 := VMRead(handle1, 0)
	if readRes1.Success != 1 {
		t.Fatalf("Read index 0 failed: %s", readRes1.String())
	}
	if readRes1.String() != "ab" {
		t.Fatalf("Expected 'ab', got '%s'", readRes1.String())
	}

	readRes2 := VMRead(handle1, 1)
	if readRes2.Success != 1 {
		t.Fatalf("Read index 1 failed: %s", readRes2.String())
	}
	if readRes2.String() != "1234567890" {
		t.Fatalf("Expected '1234567890' (truncated), got '%s'", readRes2.String())
	}

	VMClose(handle1)

	open2 := VMOpen(filename)
	if open2.Success != 1 {
		t.Fatalf("Reopen failed: %s", open2.String())
	}
	handle2, _ := strconv.Atoi(open2.String())
	defer VMClose(handle2)

	readRes3 := VMRead(handle2, 0)
	if readRes3.Success != 1 {
		t.Fatalf("Read after reopen failed: %s", readRes3.String())
	}
	if readRes3.String() != "ab" {
		t.Fatalf("After reopen: expected 'ab', got '%s'", readRes3.String())
	}

	readRes4 := VMRead(handle2, 1)
	if readRes4.Success != 1 {
		t.Fatalf("Read after reopen failed: %s", readRes4.String())
	}
	if readRes4.String() != "1234567890" {
		t.Fatalf("After reopen: expected '1234567890', got '%s'", readRes4.String())
	}
}

func TestAPIVarcharArrayReadWritePersist(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "api_varchar_persist")

	createRes := VMCreate(filename, 20, "varchar", 0)
	if createRes.Success != 1 {
		t.Fatalf("Create failed: %s", createRes.String())
	}

	open1 := VMOpen(filename)
	if open1.Success != 1 {
		t.Fatalf("Open failed: %s", open1.String())
	}
	handle1, _ := strconv.Atoi(open1.String())
	defer VMClose(handle1)

	writeRes1 := VMWrite(handle1, 0, "")
	if writeRes1.Success != 1 {
		t.Fatalf("Write empty string failed: %s", writeRes1.String())
	}

	writeRes2 := VMWrite(handle1, 3, "hello")
	if writeRes2.Success != 1 {
		t.Fatalf("Write 'hello' failed: %s", writeRes2.String())
	}

	writeRes3 := VMWrite(handle1, 7, "a much longer string")
	if writeRes3.Success != 1 {
		t.Fatalf("Write long string failed: %s", writeRes3.String())
	}

	readRes1 := VMRead(handle1, 0)
	if readRes1.Success != 1 {
		t.Fatalf("Read index 0 failed: %s", readRes1.String())
	}
	if readRes1.String() != "" {
		t.Fatalf("Expected empty string, got '%s'", readRes1.String())
	}

	readRes2 := VMRead(handle1, 3)
	if readRes2.Success != 1 {
		t.Fatalf("Read index 3 failed: %s", readRes2.String())
	}
	if readRes2.String() != "hello" {
		t.Fatalf("Expected 'hello', got '%s'", readRes2.String())
	}

	readRes3 := VMRead(handle1, 7)
	if readRes3.Success != 1 {
		t.Fatalf("Read index 7 failed: %s", readRes3.String())
	}
	if readRes3.String() != "a much longer string" {
		t.Fatalf("Expected 'a much longer string', got '%s'", readRes3.String())
	}

	VMClose(handle1)

	open2 := VMOpen(filename)
	if open2.Success != 1 {
		t.Fatalf("Reopen failed: %s", open2.String())
	}
	handle2, _ := strconv.Atoi(open2.String())
	defer VMClose(handle2)

	readRes4 := VMRead(handle2, 0)
	if readRes4.Success != 1 {
		t.Fatalf("Read after reopen failed: %s", readRes4.String())
	}
	if readRes4.String() != "" {
		t.Fatalf("After reopen: expected empty string, got '%s'", readRes4.String())
	}

	readRes5 := VMRead(handle2, 3)
	if readRes5.Success != 1 {
		t.Fatalf("Read after reopen failed: %s", readRes5.String())
	}
	if readRes5.String() != "hello" {
		t.Fatalf("After reopen: expected 'hello', got '%s'", readRes5.String())
	}

	readRes6 := VMRead(handle2, 7)
	if readRes6.Success != 1 {
		t.Fatalf("Read after reopen failed: %s", readRes6.String())
	}
	if readRes6.String() != "a much longer string" {
		t.Fatalf("After reopen: expected 'a much longer string', got '%s'", readRes6.String())
	}
}

func TestAPIHandleAfterClose(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "api_handle_close")

	createRes := VMCreate(filename, 50, "int", 0)
	if createRes.Success != 1 {
		t.Fatalf("Create failed: %s", createRes.String())
	}

	openRes := VMOpen(filename)
	if openRes.Success != 1 {
		t.Fatalf("Open failed: %s", openRes.String())
	}
	handle, _ := strconv.Atoi(openRes.String())

	writeRes := VMWrite(handle, 0, "123")
	if writeRes.Success != 1 {
		t.Fatalf("Write failed: %s", writeRes.String())
	}

	closeRes := VMClose(handle)
	if closeRes.Success != 1 {
		t.Fatalf("First close failed: %s", closeRes.String())
	}

	closeRes2 := VMClose(handle)
	if closeRes2.Success != 0 {
		t.Fatal("Expected error on second close of same handle")
	}
	if closeRes2.ErrorCode == 0 {
		t.Fatal("Expected error code to be set for second close")
	}

	readRes := VMRead(handle, 0)
	if readRes.Success != 0 {
		t.Fatal("Expected error when reading with closed handle")
	}
	if readRes.ErrorCode == 0 {
		t.Fatal("Expected error code to be set for read after close")
	}

	writeRes2 := VMWrite(handle, 1, "456")
	if writeRes2.Success != 0 {
		t.Fatal("Expected error when writing with closed handle")
	}
	if writeRes2.ErrorCode == 0 {
		t.Fatal("Expected error code to be set for write after close")
	}
}

func TestAPIStats(t *testing.T) {
	invalidStatsRes := VMStats(9999)
	if invalidStatsRes.Success != 0 {
		t.Fatal("Expected error for invalid handle in VMStats")
	}
	if invalidStatsRes.ErrorCode == 0 {
		t.Fatal("Expected error code to be set for invalid handle")
	}

	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "api_stats")

	createRes := VMCreate(filename, 50, "int", 0)
	if createRes.Success != 1 {
		t.Fatalf("Create failed: %s", createRes.String())
	}

	openRes := VMOpen(filename)
	if openRes.Success != 1 {
		t.Fatalf("Open failed: %s", openRes.String())
	}
	handle, _ := strconv.Atoi(openRes.String())
	defer VMClose(handle)

	VMWrite(handle, 0, "10")
	VMWrite(handle, 10, "20")

	statsRes := VMStats(handle)
	if statsRes.Success != 1 {
		t.Fatalf("VMStats failed: %s", statsRes.String())
	}

	stats := statsRes.String()
	if !strings.Contains(stats, "Virtual Array Stats:") {
		t.Fatal("Stats should contain 'Virtual Array Stats:'")
	}
	if !strings.Contains(stats, "Array Type:") {
		t.Fatal("Stats should contain 'Array Type:'")
	}
	if !strings.Contains(stats, "int") {
		t.Fatal("Stats should contain array type 'int'")
	}
	if !strings.Contains(stats, "Array Size:") {
		t.Fatal("Stats should contain 'Array Size:'")
	}
	if !strings.Contains(stats, "Total Pages:") {
		t.Fatal("Stats should contain 'Total Pages:'")
	}
	if !strings.Contains(stats, "Cached Pages:") {
		t.Fatal("Stats should contain 'Cached Pages:'")
	}
}
