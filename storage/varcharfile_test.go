package storage

import (
	"VirtualMemoryManagement/tests/testutils"
	"testing"
)

func TestVarcharFileCreateWriteRead(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "strings.dat")

	vf := NewVarcharFile(filename)
	if err := vf.Create(); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer vf.Close()

	offset, err := vf.GetCurrentOffset()
	if err != nil {
		t.Fatalf("GetCurrentOffset failed: %v", err)
	}

	value := "hello varchar"
	if err := vf.WriteString(offset, value); err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}

	readBack, err := vf.ReadString(offset)
	if err != nil {
		t.Fatalf("ReadString failed: %v", err)
	}
	if readBack != value {
		t.Fatalf("Expected %q, got %q", value, readBack)
	}
}

func TestVarcharFilePersistAcrossOpen(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "strings_persist.dat")

	vf := NewVarcharFile(filename)
	if err := vf.Create(); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	offset, err := vf.GetCurrentOffset()
	if err != nil {
		t.Fatalf("GetCurrentOffset failed: %v", err)
	}

	value := "persistent string"
	if err := vf.WriteString(offset, value); err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}

	if err := vf.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	vf2 := NewVarcharFile(filename)
	if err := vf2.Open(); err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer vf2.Close()

	readBack, err := vf2.ReadString(offset)
	if err != nil {
		t.Fatalf("ReadString after reopen failed: %v", err)
	}
	if readBack != value {
		t.Fatalf("Expected %q, got %q", value, readBack)
	}
}
