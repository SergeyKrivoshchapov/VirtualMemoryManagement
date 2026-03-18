package storage

import (
	"VirtualMemoryManagement/config"
	"VirtualMemoryManagement/errors"
	"VirtualMemoryManagement/types/array"
	"VirtualMemoryManagement/types/page"
	"io"
	"os"
)

const (
	SignatureSize = 2
	Signature     = "VM"
)

type PageFile struct {
	file      *os.File
	filename  string
	arrayInfo *array.Info
	binaryIO  *BinaryIO
}

var _ PageStorage = (*PageFile)(nil)

// NewPageFile creates a new PageFile instance for the given filename.
func NewPageFile(filename string) *PageFile {
	return &PageFile{
		filename: filename,
		binaryIO: NewBinaryIO(),
	}
}

// Create initializes a new page file with the specified array configuration.
func (pf *PageFile) Create(size int, typ array.Type, stringLength int) error {
	f, err := os.OpenFile(pf.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.ErrFileOperation
	}
	pf.file = f

	pf.arrayInfo = array.NewInfo(size, typ, stringLength)

	if _, err := f.WriteString(Signature); err != nil {
		return errors.ErrFileOperation
	}

	header := &Header{
		Size:         int64(size),
		Type:         typ,
		StringLength: int32(stringLength),
	}
	if err := header.WriteTo(f); err != nil {
		return errors.ErrFileOperation
	}

	for i := 0; i < pf.arrayInfo.PageCount; i++ {
		pageSize := config.TotalPageSize(pf.arrayInfo.ElementSize)
		zeroPage := make([]byte, pageSize)
		if _, err := f.Write(zeroPage); err != nil {
			return errors.ErrInsufficientDisk
		}
	}

	// Sync to ensure all data is written to disk
	if err := f.Sync(); err != nil {
		return errors.ErrFileOperation
	}

	return nil
}

// Open opens an existing page file and reads its header information.
func (pf *PageFile) Open(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return errors.ErrFileNotFound
	}
	pf.file = f
	pf.filename = filename

	sig := make([]byte, SignatureSize)
	if _, err := f.Read(sig); err != nil {
		return errors.ErrFileOperation
	}

	if string(sig) != Signature {
		return errors.ErrInvalidType
	}

	header := &Header{}
	if err := header.ReadFrom(f); err != nil {
		return errors.ErrFileOperation
	}

	pf.arrayInfo = array.NewInfo(int(header.Size), header.Type, int(header.StringLength))

	return nil
}

func (pf *PageFile) Close() error {
	if pf.file != nil {
		// Sync to ensure all pending data is written to disk
		if err := pf.file.Sync(); err != nil {
			return errors.ErrFileOperation
		}
		return pf.file.Close()
	}
	return nil
}

func (pf *PageFile) ReadPage(pageNumber int) (*page.Page, error) {
	if pageNumber < 0 || pageNumber >= pf.arrayInfo.PageCount {
		return nil, errors.ErrIndexOutOfRange
	}

	offset := pf.calculatePageOffset(pageNumber)
	if _, err := pf.file.Seek(offset, io.SeekStart); err != nil {
		return nil, errors.ErrFileOperation
	}

	bitmapData := make([]byte, config.BytesPerBitmap)
	if _, err := io.ReadFull(pf.file, bitmapData); err != nil {
		return nil, errors.ErrFileOperation
	}

	pageData := make([]byte, config.PageDataSize(pf.arrayInfo.ElementSize))
	if _, err := io.ReadFull(pf.file, pageData); err != nil {
		return nil, errors.ErrFileOperation
	}

	p, err := page.NewWithData(pageNumber, pf.arrayInfo.ElementSize, bitmapData, pageData)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (pf *PageFile) WritePage(p *page.Page) error {
	if p.AbsoluteNumber < 0 || p.AbsoluteNumber >= pf.arrayInfo.PageCount {
		return errors.ErrIndexOutOfRange
	}

	offset := pf.calculatePageOffset(p.AbsoluteNumber)
	if _, err := pf.file.Seek(offset, io.SeekStart); err != nil {
		return errors.ErrFileOperation
	}

	if _, err := pf.file.Write(p.Bitmap().Bytes()); err != nil {
		return errors.ErrFileOperation
	}

	if _, err := pf.file.Write(p.Data()); err != nil {
		return errors.ErrFileOperation
	}

	// Sync to ensure data is written to disk
	if err := pf.file.Sync(); err != nil {
		return errors.ErrFileOperation
	}

	return nil
}

func (pf *PageFile) ArrayInfo() *array.Info {
	return pf.arrayInfo
}

func (pf *PageFile) calculatePageOffset(pageNumber int) int64 {
	headerSize := int64(SignatureSize) + int64((&Header{}).Size_())
	pageSize := int64(config.TotalPageSize(pf.arrayInfo.ElementSize))
	return headerSize + int64(pageNumber)*pageSize
}
