package storage

import (
	"bytes"
	"testing"
)

func TestBinaryIOWriteReadInt32(t *testing.T) {
	bio := NewBinaryIO()
	buf := new(bytes.Buffer)

	testValue := int32(42)
	err := bio.WriteInt32(buf, testValue)
	if err != nil {
		t.Fatalf("Failed to write int32: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	readValue, err := bio.ReadInt32(reader)
	if err != nil {
		t.Fatalf("Failed to read int32: %v", err)
	}

	if readValue != testValue {
		t.Fatalf("Expected %d, got %d", testValue, readValue)
	}
}

func TestBinaryIOWriteReadInt64(t *testing.T) {
	bio := NewBinaryIO()
	buf := new(bytes.Buffer)

	testValue := int64(1234567890)
	err := bio.WriteInt64(buf, testValue)
	if err != nil {
		t.Fatalf("Failed to write int64: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	readValue, err := bio.ReadInt64(reader)
	if err != nil {
		t.Fatalf("Failed to read int64: %v", err)
	}

	if readValue != testValue {
		t.Fatalf("Expected %d, got %d", testValue, readValue)
	}
}

func TestBinaryIOWriteReadByte(t *testing.T) {
	bio := NewBinaryIO()
	buf := new(bytes.Buffer)

	testValue := byte(255)
	err := bio.WriteByte(buf, testValue)
	if err != nil {
		t.Fatalf("Failed to write byte: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	readValue, err := bio.ReadByte(reader)
	if err != nil {
		t.Fatalf("Failed to read byte: %v", err)
	}

	if readValue != testValue {
		t.Fatalf("Expected %d, got %d", testValue, readValue)
	}
}

func TestBinaryIOWriteReadBytes(t *testing.T) {
	bio := NewBinaryIO()
	buf := new(bytes.Buffer)

	testData := []byte{1, 2, 3, 4, 5}
	err := bio.WriteBytes(buf, testData)
	if err != nil {
		t.Fatalf("Failed to write bytes: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	readData := make([]byte, len(testData))
	err = bio.ReadBytes(reader, readData)
	if err != nil {
		t.Fatalf("Failed to read bytes: %v", err)
	}

	if !bytes.Equal(readData, testData) {
		t.Fatalf("Data mismatch: expected %v, got %v", testData, readData)
	}
}

func TestBinaryIOMultipleValues(t *testing.T) {
	bio := NewBinaryIO()
	buf := new(bytes.Buffer)

	val1 := int32(100)
	val2 := int64(200000)
	val3 := byte(50)

	bio.WriteInt32(buf, val1)
	bio.WriteInt64(buf, val2)
	bio.WriteByte(buf, val3)

	reader := bytes.NewReader(buf.Bytes())

	r1, _ := bio.ReadInt32(reader)
	r2, _ := bio.ReadInt64(reader)
	r3, _ := bio.ReadByte(reader)

	if r1 != val1 || r2 != val2 || r3 != val3 {
		t.Fatal("Values mismatch after multiple writes and reads")
	}
}

func TestBinaryIOLargeInt32Values(t *testing.T) {
	bio := NewBinaryIO()

	testCases := []int32{0, 1, -1, 2147483647, -2147483648}

	for _, val := range testCases {
		buf := new(bytes.Buffer)
		bio.WriteInt32(buf, val)

		reader := bytes.NewReader(buf.Bytes())
		readVal, err := bio.ReadInt32(reader)
		if err != nil {
			t.Fatalf("Failed for value %d: %v", val, err)
		}
		if readVal != val {
			t.Fatalf("Expected %d, got %d", val, readVal)
		}
	}
}

func TestBinaryIOLargeInt64Values(t *testing.T) {
	bio := NewBinaryIO()

	testCases := []int64{0, 1, -1, 9223372036854775807, -9223372036854775808}

	for _, val := range testCases {
		buf := new(bytes.Buffer)
		bio.WriteInt64(buf, val)

		reader := bytes.NewReader(buf.Bytes())
		readVal, err := bio.ReadInt64(reader)
		if err != nil {
			t.Fatalf("Failed for value %d: %v", val, err)
		}
		if readVal != val {
			t.Fatalf("Expected %d, got %d", val, readVal)
		}
	}
}

func TestBinaryIOZeroBytes(t *testing.T) {
	bio := NewBinaryIO()
	buf := new(bytes.Buffer)

	testData := make([]byte, 100)
	for i := range testData {
		testData[i] = 0
	}

	bio.WriteBytes(buf, testData)
	reader := bytes.NewReader(buf.Bytes())

	readData := make([]byte, len(testData))
	bio.ReadBytes(reader, readData)

	if !bytes.Equal(readData, testData) {
		t.Fatal("Zero bytes not preserved")
	}
}

func TestBinaryIOStructToBytes(t *testing.T) {
	bio := NewBinaryIO()

	type TestStruct struct {
		A int32
		B int32
	}

	ts := TestStruct{A: 100, B: 200}
	data := bio.StructToBytes(ts)

	if len(data) == 0 {
		t.Fatal("StructToBytes returned empty data")
	}
}

func TestBinaryIOMaxByteValue(t *testing.T) {
	bio := NewBinaryIO()

	for val := byte(0); val < 255; val++ {
		buf := new(bytes.Buffer)
		bio.WriteByte(buf, val)

		reader := bytes.NewReader(buf.Bytes())
		readVal, _ := bio.ReadByte(reader)
		if readVal != val {
			t.Fatalf("Byte value mismatch for %d", val)
		}
	}
}
