// Package bitmap provides a bitmap implementation for tracking allocated slots.
package bitmap

import (
	"VirtualMemoryManagement/config"
	"VirtualMemoryManagement/errors"
	"fmt"
)

type BitMap struct {
	bits [config.BytesPerBitmap]byte
}

var ErrBitPositionOutOfRange = errors.NewError(errors.ErrCodeIndexOutOfRange, "bit position out of range")

func New() *BitMap {
	return &BitMap{}
}

func (bm *BitMap) Set(pos int) error {
	if pos < 0 || pos >= config.BitsPerPage {
		return errors.NewErrorWithWrapped(errors.ErrCodeIndexOutOfRange, 
			fmt.Sprintf("bit position %d out of range [0,%d]", pos, config.BitsPerPage-1), 
			ErrBitPositionOutOfRange)
	}
	bytePos := pos / 8
	bitPos := uint(pos % 8)
	bm.bits[bytePos] |= 1 << bitPos
	return nil
}

func (bm *BitMap) IsSet(pos int) (bool, error) {
	if pos < 0 || pos >= config.BitsPerPage {
		return false, errors.NewErrorWithWrapped(errors.ErrCodeIndexOutOfRange,
			fmt.Sprintf("bit position %d out of range [0,%d]", pos, config.BitsPerPage-1),
			ErrBitPositionOutOfRange)
	}
	bytePos := pos / 8
	bitPos := uint(pos % 8)
	return bm.bits[bytePos]&(1<<bitPos) != 0, nil
}

func (bm *BitMap) Bytes() []byte {
	return bm.bits[:]
}

func (bm *BitMap) FromBytes(data []byte) error {
	if len(data) != config.BytesPerBitmap {
		return errors.NewError(errors.ErrCodeFileOperation, 
			fmt.Sprintf("bitmap requires %d bytes, got %d", config.BytesPerBitmap, len(data)))
	}
	copy(bm.bits[:], data)
	return nil
}
