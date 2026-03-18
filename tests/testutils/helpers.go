package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

// KEEP_TEST_FILES_ENV - если установить переменную окружения, тестовые файлы не будут удаляться
// Используется для отладки: export KEEP_TEST_FILES=1
const KEEP_TEST_FILES_ENV = "KEEP_TEST_FILES"

func TempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "vmm_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

func TempFile(t *testing.T, dir, pattern string) string {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func CleanupDir(t *testing.T, dir string) {
	// Если установлена переменная окружения KEEP_TEST_FILES, не удаляем файлы
	if os.Getenv(KEEP_TEST_FILES_ENV) != "" {
		t.Logf("Keeping test files at: %s", dir)
		return
	}

	if err := os.RemoveAll(dir); err != nil {
		t.Logf("Warning: failed to cleanup dir %s: %v", dir, err)
	}
}

func RemoveFile(t *testing.T, path string) {
	// Если установлена переменная окружения KEEP_TEST_FILES, не удаляем файлы
	if os.Getenv(KEEP_TEST_FILES_ENV) != "" {
		return
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Logf("Warning: failed to remove file %s: %v", path, err)
	}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return info.Size()
}

func ReadFileBytes(t *testing.T, path string, offset int64, count int) []byte {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	f.Seek(offset, 0)
	data := make([]byte, count)
	n, err := f.Read(data)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	return data[:n]
}

func TempFilePath(dir, name string) string {
	return filepath.Join(dir, name)
}
