package storage

import (
	"VirtualMemoryManagement/config"
	"VirtualMemoryManagement/tests/testutils"
	"VirtualMemoryManagement/types/array"
	"os"
	"testing"
)

func TestPageFileCreateAndOpenInt(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)
	filename := testutils.TempFilePath(dir, "pf_int.vmm")
	pf := NewPageFile(filename)
	size := 1000
	err := pf.Create(size, array.TypeInt, 0)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := pf.Close(); err != nil {
		t.Fatalf("Close after create failed: %v", err)
	}
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open created file: %v", err)
	}
	defer f.Close()

	sig := make([]byte, SignatureSize)
	if _, err := f.Read(sig); err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}
	if string(sig) != Signature {
		t.Fatalf("Expected signature %q, got %q", Signature, string(sig))
	}
	header := &Header{}
	if err := header.ReadFrom(f); err != nil {
		t.Fatalf("Failed to read header: %v", err)
	}
	if header.Size != int64(size) {
		t.Fatalf("Expected size %d, got %d", size, header.Size)
	}
	if header.Type != array.TypeInt {
		t.Fatalf("Expected type %d, got %d", array.TypeInt, header.Type)
	}
	if header.StringLength != 0 {
		t.Fatalf("Expected length %d, got %d", 0, header.StringLength)
	}

	info := array.NewInfo(size, array.TypeInt, 0)
	expectedDataSizePerPage := config.PageDataSize(info.ElementSize)
	expectedTotalPageSize := config.TotalPageSize(info.ElementSize)
	stat, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed stat: %v", err)
	}
	actualSize := stat.Size()
	headerSize := int64(SignatureSize) + int64((&Header{}).Size_())
	expectedSize := headerSize + int64(info.PageCount)*int64(expectedTotalPageSize)
	if actualSize != expectedSize {
		t.Fatalf("Expected size %d, got %d", expectedSize, actualSize)
	}
	if expectedDataSizePerPage != config.PageDataSize(info.ElementSize) {
		t.Fatalf("Unexpected page data size")
	}
}

func TestPageFileReadWritePagePersists(t *testing.T) {
	dir := testutils.TempDir(t)
	defer testutils.CleanupDir(t, dir)

	filename := testutils.TempFilePath(dir, "pf_rw.vmm")
	size := 128

	pf := NewPageFile(filename)
	if err := pf.Create(size, array.TypeInt, 0); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	info := pf.ArrayInfo()
	if info.PageCount != 1 {
		t.Fatalf("Expected 1 page, got %d", info.PageCount)
	}

	p, err := pf.ReadPage(0)
	if err != nil {
		t.Fatalf("ReadPage failed: %v", err)
	}

	if err := p.SetBit(0); err != nil {
		t.Fatalf("SetBit failed: %v", err)
	}

	p.Data()[0] = 0x12
	p.Data()[1] = 0x34
	p.Data()[2] = 0x56
	p.Data()[3] = 0x78

	if err := pf.WritePage(p); err != nil {
		t.Fatalf("WritePage failed: %v", err)
	}

	if err := pf.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	pf2 := NewPageFile(filename)
	if err := pf2.Open(filename); err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer pf2.Close()

	p2, err := pf2.ReadPage(0)
	if err != nil {
		t.Fatalf("ReadPage after reopen failed: %v", err)
	}

	isSet, err := p2.IsBitSet(0)
	if err != nil {
		t.Fatalf("IsBitSet failed: %v", err)
	}
	if !isSet {
		t.Fatalf("Expected bit 0 to be set after reload")
	}

	data2 := p2.Data()
	if data2[0] != 0x12 || data2[1] != 0x34 || data2[2] != 0x56 || data2[3] != 0x78 {
		t.Fatalf("Data mismatch after reload, got %v", data2[:4])
	}
}
