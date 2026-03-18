package api

import (
	"VirtualMemoryManagement/tests/testutils"
	"os"
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

// TestVMWriteVarchar проверяет запись и чтение строк переменной длины
func TestVMWriteVarchar(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_varchar_write")

	createResult := VMCreate(filename, 100, "varchar", 100)
	if createResult.Success != 1 {
		t.Fatalf("Create failed: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Open failed: %s", openResult.String())
	}

	handle, err := strconv.Atoi(openResult.String())
	if err != nil {
		t.Fatalf("Failed to parse handle: %v", err)
	}

	testStrings := []struct {
		index int
		value string
	}{
		{0, "Hello"},
		{1, "World"},
		{2, "Virtual Memory"},
		{5, "Test String"},
		{10, "abc123"},
		{50, "Longer string with spaces"},
	}

	// Test Case 1: Записать строки разной длины
	for _, ts := range testStrings {
		writeResult := VMWrite(handle, ts.index, ts.value)
		if writeResult.Success != 1 {
			t.Fatalf("Write at index %d failed: %s", ts.index, writeResult.String())
		}
	}

	// Test Case 2: Прочитать и проверить записанные строки
	for _, ts := range testStrings {
		readResult := VMRead(handle, ts.index)
		if readResult.Success != 1 {
			t.Fatalf("Read at index %d failed: %s", ts.index, readResult.String())
		}

		actualValue := readResult.String()
		if actualValue != ts.value {
			t.Errorf("Index %d: expected %q, got %q", ts.index, ts.value, actualValue)
		}
	}

	// Test Case 3: Записать пустую строку
	emptyWriteResult := VMWrite(handle, 20, "")
	if emptyWriteResult.Success != 1 {
		t.Fatalf("Write empty string failed: %s", emptyWriteResult.String())
	}

	emptyReadResult := VMRead(handle, 20)
	if emptyReadResult.Success != 1 {
		t.Fatalf("Read empty string failed: %s", emptyReadResult.String())
	}

	if emptyReadResult.String() != "" {
		t.Fatalf("Expected empty string, got %q", emptyReadResult.String())
	}

	// Test Case 4: Перезаписать существующее значение
	originalValue := "Original"
	newValue := "Updated"

	VMWrite(handle, 30, originalValue)
	readOrig := VMRead(handle, 30)
	if readOrig.String() != originalValue {
		t.Fatalf("Original value mismatch: expected %q, got %q", originalValue, readOrig.String())
	}

	// Перезаписать
	updateResult := VMWrite(handle, 30, newValue)
	if updateResult.Success != 1 {
		t.Fatalf("Update write failed: %s", updateResult.String())
	}

	readUpdated := VMRead(handle, 30)
	if readUpdated.String() != newValue {
		t.Fatalf("Updated value mismatch: expected %q, got %q", newValue, readUpdated.String())
	}

	// Test Case 5: Запись и чтение специальных символов
	specialChars := []string{
		"123!@#$%^&*()",
		"привет мир",
		"Line1\nLine2",
		"Tab\tSeparated",
		"Quote: \"test\"",
		"   spaces   ",
	}

	for i, special := range specialChars {
		idx := 40 + i
		writeResult := VMWrite(handle, idx, special)
		if writeResult.Success != 1 {
			t.Fatalf("Write special chars at index %d failed: %s", idx, writeResult.String())
		}

		readResult := VMRead(handle, idx)
		if readResult.Success != 1 {
			t.Fatalf("Read special chars at index %d failed: %s", idx, readResult.String())
		}

		if readResult.String() != special {
			t.Errorf("Special chars mismatch at index %d: expected %q, got %q",
				idx, special, readResult.String())
		}
	}

	// Test Case 6: Получить статистику
	statsResult := VMStats(handle)
	if statsResult.Success != 1 {
		t.Fatalf("VMStats failed: %s", statsResult.String())
	}

	stats := statsResult.String()
	if !strings.Contains(stats, "varchar") && !strings.Contains(stats, "V") {
		t.Fatal("Stats should contain 'varchar' or 'V'")
	}

	// Закрыть файл
	closeResult := VMClose(handle)
	if closeResult.Success != 1 {
		t.Fatalf("Close failed: %s", closeResult.String())
	}

	// Test Case 7: Переоткрыть файл и проверить что данные сохранились
	reopenResult := VMOpen(filename)
	if reopenResult.Success != 1 {
		t.Fatalf("Reopen failed: %s", reopenResult.String())
	}

	handle2, err := strconv.Atoi(reopenResult.String())
	if err != nil {
		t.Fatalf("Failed to parse handle: %v", err)
	}

	// Проверить что данные по-прежнему там
	for _, ts := range testStrings {
		readResult := VMRead(handle2, ts.index)
		if readResult.Success != 1 {
			t.Fatalf("Persistence read at index %d failed: %s", ts.index, readResult.String())
		}

		actualValue := readResult.String()
		if actualValue != ts.value {
			t.Errorf("Persistence check at index %d: expected %q, got %q",
				ts.index, ts.value, actualValue)
		}
	}

	VMClose(handle2)
}

// TestVMWriteVarcharLongStrings проверяет запись длинных строк
func TestVMWriteVarcharLongStrings(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_varchar_long")

	// Создать varchar массив с максимальной длиной 200 символов
	createResult := VMCreate(filename, 50, "varchar", 200)
	if createResult.Success != 1 {
		t.Fatalf("Create failed: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Open failed: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())

	// Создать длинные строки
	longStrings := []string{
		"This is a longer string with multiple words and characters.",
		"The quick brown fox jumps over the lazy dog. " +
			"This sentence contains all letters of the English alphabet.",
		"String with numbers: 0123456789 and special: !@#$%^&*()",
		"Very long string: " + strings.Repeat("a", 100),
		"Another test: " + strings.Repeat("xyz", 30),
	}

	// Записать длинные строки
	for i, str := range longStrings {
		writeResult := VMWrite(handle, i, str)
		if writeResult.Success != 1 {
			t.Fatalf("Write long string %d failed: %s", i, writeResult.String())
		}
	}

	// Прочитать и проверить
	for i, expectedStr := range longStrings {
		readResult := VMRead(handle, i)
		if readResult.Success != 1 {
			t.Fatalf("Read long string %d failed: %s", i, readResult.String())
		}

		actualStr := readResult.String()
		if actualStr != expectedStr {
			t.Errorf("Long string %d mismatch:\nExpected: %q\nGot: %q",
				i, expectedStr, actualStr)
		}
	}

	VMClose(handle)
}

// TestVMWriteVarcharIndexOutOfRange проверяет обработку ошибок
func TestVMWriteVarcharIndexOutOfRange(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_varchar_bounds")

	// Создать с размером 10
	createResult := VMCreate(filename, 10, "varchar", 100)
	if createResult.Success != 1 {
		t.Fatalf("Create failed: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Open failed: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())

	// Попробовать записать вне диапазона (индекс 10 для size 10)
	writeResult := VMWrite(handle, 10, "out of range")
	if writeResult.Success != 0 {
		t.Fatal("Expected error for out of range index")
	}

	// Попробовать записать отрицательный индекс
	writeResult = VMWrite(handle, -1, "negative")
	if writeResult.Success != 0 {
		t.Fatal("Expected error for negative index")
	}

	// Проверить что валидные индексы работают (0-9)
	for i := 0; i < 10; i++ {
		writeResult := VMWrite(handle, i, "valid")
		if writeResult.Success != 1 {
			t.Fatalf("Write at valid index %d failed", i)
		}
	}

	VMClose(handle)
}

// TestVMHelpBasic проверяет что VMHelp пишет текст в файл
func TestVMHelpBasic(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	helpFilePath := testutils.TempFilePath(dir, "help")

	helpText := "Virtual Memory Manager Help\n" +
		"Commands:\n" +
		"  Create <filename> <size> <type> - Create array\n" +
		"  Open <filename> - Open existing array\n" +
		"  Close <handle> - Close array\n" +
		"  Read <handle> <index> - Read element\n" +
		"  Write <handle> <index> <value> - Write element\n"

	result := VMHelp(helpFilePath, helpText)

	// Проверить результат
	if result.Success != 1 {
		t.Fatalf("VMHelp failed: %s", result.String())
	}

	// Проверить что файл был создан
	if !testutils.FileExists(helpFilePath) {
		t.Fatalf("Help file was not created at %s", helpFilePath)
	}

	// Прочитать содержимое файла
	content, err := os.ReadFile(helpFilePath)
	if err != nil {
		t.Fatalf("Failed to read help file: %v", err)
	}

	// Проверить что содержимое совпадает
	if string(content) != helpText {
		t.Fatalf("Help text mismatch.\nExpected:\n%s\n\nGot:\n%s", helpText, string(content))
	}
}

// TestVMHelpEmptyText проверяет запись пустого текста
func TestVMHelpEmptyText(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	helpFilePath := testutils.TempFilePath(dir, "help_empty")

	result := VMHelp(helpFilePath, "")

	if result.Success != 1 {
		t.Fatalf("VMHelp with empty text failed: %s", result.String())
	}

	if !testutils.FileExists(helpFilePath) {
		t.Fatalf("Help file was not created")
	}

	content, _ := os.ReadFile(helpFilePath)
	if len(content) != 0 {
		t.Fatalf("Expected empty file, got %d bytes", len(content))
	}
}

// TestVMHelpMultiline проверяет запись многострочного текста
func TestVMHelpMultiline(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	helpFilePath := testutils.TempFilePath(dir, "help_multiline")

	helpText := "Line 1\nLine 2\nLine 3\n\nLine 5 (after empty line)"

	result := VMHelp(helpFilePath, helpText)

	if result.Success != 1 {
		t.Fatalf("VMHelp with multiline text failed: %s", result.String())
	}

	content, _ := os.ReadFile(helpFilePath)
	if string(content) != helpText {
		t.Fatalf("Multiline text mismatch")
	}
}

// TestVMHelpSpecialCharacters проверяет текст со спецсимволами
func TestVMHelpSpecialCharacters(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	helpFilePath := testutils.TempFilePath(dir, "help_special")

	helpText := "Special: !@#$%^&*() []{} <>\n" +
		"Unicode: привет мир 你好\n" +
		"Quotes: \"double\" 'single'\n" +
		"Tabs:\tand\ttabs\n"

	result := VMHelp(helpFilePath, helpText)

	if result.Success != 1 {
		t.Fatalf("VMHelp with special chars failed: %s", result.String())
	}

	content, _ := os.ReadFile(helpFilePath)
	if string(content) != helpText {
		t.Fatalf("Special characters text mismatch")
	}
}

// TestVMHelpLargeText проверяет запись большого текста
func TestVMHelpLargeText(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	helpFilePath := testutils.TempFilePath(dir, "help_large")

	// Создать большой текст
	helpText := strings.Repeat("This is help text. ", 1000) // ~20KB текста

	result := VMHelp(helpFilePath, helpText)

	if result.Success != 1 {
		t.Fatalf("VMHelp with large text failed: %s", result.String())
	}

	content, _ := os.ReadFile(helpFilePath)
	if string(content) != helpText {
		t.Fatalf("Large text mismatch")
	}
}

// TestVMHelpFileOverwrite проверяет перезапись существующего файла
func TestVMHelpFileOverwrite(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	helpFilePath := testutils.TempFilePath(dir, "help_overwrite")

	// Первая запись
	firstText := "First help text"
	result1 := VMHelp(helpFilePath, firstText)
	if result1.Success != 1 {
		t.Fatalf("First VMHelp failed: %s", result1.String())
	}

	content1, _ := os.ReadFile(helpFilePath)
	if string(content1) != firstText {
		t.Fatalf("First write mismatch")
	}

	// Вторая запись (перезапись)
	secondText := "Second help text - different content"
	result2 := VMHelp(helpFilePath, secondText)
	if result2.Success != 1 {
		t.Fatalf("Second VMHelp failed: %s", result2.String())
	}

	content2, _ := os.ReadFile(helpFilePath)
	if string(content2) != secondText {
		t.Fatalf("Second write mismatch")
	}

	// Убедиться что первый текст больше не там
	if string(content2) == firstText {
		t.Fatal("File was not overwritten")
	}
}

// TestVMHelpInvalidPath проверяет обработку ошибок при неверном пути
func TestVMHelpInvalidPath(t *testing.T) {
	// Пытаться писать в несуществующую директорию
	invalidPath := "/nonexistent/directory/that/does/not/exist/help.txt"

	result := VMHelp(invalidPath, "test text")

	if result.Success != 0 {
		t.Fatal("Expected error for invalid path")
	}

	if result.ErrorCode == 0 {
		t.Fatal("Expected error code to be set")
	}

	if !strings.Contains(result.String(), "failed to write help to file") {
		t.Fatalf("Expected error message about file write, got: %s", result.String())
	}
}

// TestVMHelpWithBinaryData проверяет текст с бинарными данными
func TestVMHelpWithBinaryData(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	helpFilePath := testutils.TempFilePath(dir, "help_binary")

	// Текст с нулевыми байтами
	helpText := "Text with null\x00byte and special\x01\x02chars\xFF"

	result := VMHelp(helpFilePath, helpText)

	if result.Success != 1 {
		t.Fatalf("VMHelp with binary data failed: %s", result.String())
	}

	content, _ := os.ReadFile(helpFilePath)
	if string(content) != helpText {
		t.Fatalf("Binary data text mismatch")
	}
}

func TestVMWriteCharSimple(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_char_simple")

	t.Logf("Creating CHAR array at: %s", filename)

	createResult := VMCreate(filename, 20, "char", 30)
	if createResult.Success != 1 {
		t.Fatalf("Create failed: %s", createResult.String())
	}
	t.Logf("CHAR array created successfully")

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Open failed: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())
	t.Logf("File opened with handle: %d", handle)

	testData := []struct {
		index int
		value string
	}{
		{0, "Hello"},
		{1, "World"},
		{2, "Test String"},
		{3, "1234567890"},
		{4, "Short"},
		{5, "This is a longer test string"},
	}

	t.Logf("\nWriting test data")
	for _, td := range testData {
		writeResult := VMWrite(handle, td.index, td.value)
		if writeResult.Success != 1 {
			t.Fatalf("Write at index %d failed: %s", td.index, writeResult.String())
		}
		t.Logf("  [%d] = %q", td.index, td.value)
	}

	t.Logf("\nReading test data")
	for _, td := range testData {
		readResult := VMRead(handle, td.index)
		if readResult.Success != 1 {
			t.Fatalf("Read at index %d failed: %s", td.index, readResult.String())
		}
		actualValue := readResult.String()
		t.Logf("  [%d] = %q", td.index, actualValue)

		if actualValue != td.value {
			t.Errorf("Mismatch at index %d: expected %q, got %q", td.index, td.value, actualValue)
		}
	}

	closeResult := VMClose(handle)
	if closeResult.Success != 1 {
		t.Fatalf("Close failed: %s", closeResult.String())
	}
	t.Logf("\nFile closed and saved")
	t.Logf("\nFile location: %s", filename)
	t.Logf("You can now inspect the .vm file manually")
}

func TestVMWriteIntOverflow(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "test_int_overflow")

	createResult := VMCreate(filename, 10, "int", 0)
	if createResult.Success != 1 {
		t.Fatalf("Create failed: %s", createResult.String())
	}

	openResult := VMOpen(filename)
	if openResult.Success != 1 {
		t.Fatalf("Open failed: %s", openResult.String())
	}

	handle, _ := strconv.Atoi(openResult.String())

	testCases := []struct {
		value string
		name  string
	}{
		{"99999999999999999999", "overflow positive"},
		{"-99999999999999999999", "overflow negative"},
		{"abc", "non-numeric"},
		{"12.34", "float"},
		{"", "empty string"},
		{"12a34", "mixed"},
	}

	for _, tc := range testCases {
		writeResult := VMWrite(handle, 0, tc.value)
		if writeResult.Success != 0 {
			t.Errorf("%s: expected error, got success", tc.name)
		}
		if writeResult.ErrorCode == 0 {
			t.Errorf("%s: expected error code, got 0", tc.name)
		}
	}

	validResult := VMWrite(handle, 0, "42")
	if validResult.Success != 1 {
		t.Fatalf("Valid int write failed: %s", validResult.String())
	}

	readResult := VMRead(handle, 0)
	if readResult.String() != "42" {
		t.Fatalf("Expected 42, got %s", readResult.String())
	}

	VMClose(handle)
}
