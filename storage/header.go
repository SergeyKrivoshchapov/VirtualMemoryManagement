package storage

import (
	"VirtualMemoryManagement/types/array"
	"bytes"
	"encoding/binary"
	"io"
)

type Header struct {
	Size         int64
	Type         array.Type
	StringLength int32
}

func (h *Header) WriteTo(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, h.Size); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, h.Type); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, h.StringLength); err != nil {
		return err
	}
	return nil
}

func (h *Header) ReadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &h.Size); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Type); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.StringLength); err != nil {
		return err
	}
	return nil
}

func (h *Header) Size_() int {
	buf := new(bytes.Buffer)
	h.WriteTo(buf)
	return buf.Len()
}

