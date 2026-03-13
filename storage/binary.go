package storage

import (
	"bytes"
	"encoding/binary"
	"io"
)

type BinaryIO struct{}

func NewBinaryIO() *BinaryIO {
	return &BinaryIO{}
}

func (b *BinaryIO) WriteBytes(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}

func (b *BinaryIO) WriteInt32(w io.Writer, value int32) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func (b *BinaryIO) WriteInt64(w io.Writer, value int64) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func (b *BinaryIO) WriteByte(w io.Writer, value byte) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func (b *BinaryIO) ReadBytes(r io.Reader, data []byte) error {
	_, err := io.ReadFull(r, data)
	return err
}

func (b *BinaryIO) ReadInt32(r io.Reader) (int32, error) {
	var val int32
	err := binary.Read(r, binary.LittleEndian, &val)
	return val, err
}

func (b *BinaryIO) ReadInt64(r io.Reader) (int64, error) {
	var val int64
	err := binary.Read(r, binary.LittleEndian, &val)
	return val, err
}

func (b *BinaryIO) ReadByte(r io.Reader) (byte, error) {
	var val byte
	err := binary.Read(r, binary.LittleEndian, &val)
	return val, err
}

func (b *BinaryIO) StructToBytes(v interface{}) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes()
}

