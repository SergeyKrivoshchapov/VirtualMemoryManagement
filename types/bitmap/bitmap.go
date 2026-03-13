package bitmap

import (
	"VirtualMemoryManagement/config"
	"fmt"
)

type BitMap struct {
	bits [config.BytesPerBitmap]byte
}

func New() *BitMap {
	return &BitMap{}
}

func (bm *BitMap) Set(pos int) error {
	if pos < 0 || pos >= config.BitsPerPage {
		return fmt.Errorf("bit position %d out of range[0,%d]", pos, config.BitsPerPage-1)
	}
	bytePos := pos / 8
	bitPos := uint(pos % 8)
	bm.bits[bytePos] |= 1 << bitPos
	return nil
}

func (bm *BitMap) IsSet(pos int) (bool, error) {
	if pos < 0 || pos >= config.BitsPerPage {
		return false, fmt.Errorf("bit position %d out of range[0,%d]", pos, config.BitsPerPage-1)
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
		return fmt.Errorf("bitmap requires %d bytes, got %d", config.BytesPerBitmap, len(data))
	}
	copy(bm.bits[:], data)
	return nil
}
