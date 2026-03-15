package storage

import (
	"VirtualMemoryManagement/errors"
	"encoding/binary"
	"io"
	"os"
)

type VarcharFile struct {
	file     *os.File
	filename string
	binaryIO *BinaryIO
	offset   int64
}

var _ VarcharStorage = (*VarcharFile)(nil)

func NewVarcharFile(filename string) *VarcharFile {
	return &VarcharFile{
		filename: filename,
		binaryIO: NewBinaryIO(),
	}
}

func (vf *VarcharFile) Create() error {
	f, err := os.Create(vf.filename)
	if err != nil {
		return errors.ErrFileOperation
	}
	vf.file = f
	vf.offset = 0
	if _, err := vf.file.Write([]byte{0}); err != nil {
		return errors.ErrFileOperation
	}
	vf.offset = 1
	return nil
}

func (vf *VarcharFile) Open() error {
	f, err := os.OpenFile(vf.filename, os.O_RDWR, 0644)
	if err != nil {
		return errors.ErrFileNotFound
	}
	vf.file = f
	info, err := f.Stat()
	if err == nil {
		vf.offset = info.Size()
		if vf.offset == 0 {
			vf.offset = 1
		}
	} else {
		vf.offset = 1
	}
	return nil
}

func (vf *VarcharFile) Close() error {
	if vf.file != nil {
		return vf.file.Close()
	}
	return nil
}

func (vf *VarcharFile) WriteString(offset int64, s string) error {
	if _, err := vf.file.Seek(offset, io.SeekStart); err != nil {
		return errors.ErrFileOperation
	}

	if err := binary.Write(vf.file, binary.LittleEndian, int32(len(s))); err != nil {
		return errors.ErrFileOperation
	}

	if _, err := vf.file.WriteString(s); err != nil {
		return errors.ErrFileOperation
	}

	if err := vf.file.Sync(); err != nil {
		return errors.ErrFileOperation
	}

	vf.offset = offset + 4 + int64(len(s))
	return nil
}

func (vf *VarcharFile) ReadString(offset int64) (string, error) {
	if _, err := vf.file.Seek(offset, io.SeekStart); err != nil {
		return "", errors.ErrFileOperation
	}

	var length int32
	if err := binary.Read(vf.file, binary.LittleEndian, &length); err != nil {
		return "", errors.ErrFileOperation
	}

	if length < 0 || length > 1000000 {
		return "", errors.NewError(errors.ErrCodeFileOperation, "Invalid string length")
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(vf.file, data); err != nil {
		return "", errors.ErrFileOperation
	}

	return string(data), nil
}

func (vf *VarcharFile) GetCurrentOffset() (int64, error) {
	return vf.offset, nil
}
