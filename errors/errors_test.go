package errors

import (
	"testing"
)

func TestErrorFileNotFound(t *testing.T) {
	err := ErrFileNotFound
	if err == nil {
		t.Fatal("ErrFileNotFound should not be nil")
	}
	if err.Code != ErrCodeFileNotFound {
		t.Fatalf("Expected code %d, got %d", ErrCodeFileNotFound, err.Code)
	}
}

func TestErrorOutOfMemory(t *testing.T) {
	err := ErrOutOfMemory
	if err.Code != ErrCodeOutOfMemory {
		t.Fatalf("Expected code %d, got %d", ErrCodeOutOfMemory, err.Code)
	}
}

func TestErrorIndexOutOfRange(t *testing.T) {
	err := ErrIndexOutOfRange
	if err.Code != ErrCodeIndexOutOfRange {
		t.Fatalf("Expected code %d, got %d", ErrCodeIndexOutOfRange, err.Code)
	}
}

func TestErrorFileOperation(t *testing.T) {
	err := ErrFileOperation
	if err.Code != ErrCodeFileOperation {
		t.Fatalf("Expected code %d, got %d", ErrCodeFileOperation, err.Code)
	}
}

func TestErrorInvalidType(t *testing.T) {
	err := ErrInvalidType
	if err.Code != ErrCodeInvalidType {
		t.Fatalf("Expected code %d, got %d", ErrCodeInvalidType, err.Code)
	}
}

func TestErrorInsufficientDisk(t *testing.T) {
	err := ErrInsufficientDisk
	if err.Code != ErrCodeInsufficientDisk {
		t.Fatalf("Expected code %d, got %d", ErrCodeInsufficientDisk, err.Code)
	}
}

func TestErrorInvalidHandle(t *testing.T) {
	err := ErrInvalidHandle
	if err.Code != ErrCodeInvalidHandle {
		t.Fatalf("Expected code %d, got %d", ErrCodeInvalidHandle, err.Code)
	}
}

func TestErrorPageNotFound(t *testing.T) {
	err := ErrPageNotFound
	if err.Code != ErrCodePageNotFound {
		t.Fatalf("Expected code %d, got %d", ErrCodePageNotFound, err.Code)
	}
}

func TestNewError(t *testing.T) {
	err := NewError(42, "test error")
	if err.Code != 42 {
		t.Fatalf("Expected code 42, got %d", err.Code)
	}
	if err.Message != "test error" {
		t.Fatalf("Expected message 'test error', got '%s'", err.Message)
	}
}

func TestNewErrorWithWrapped(t *testing.T) {
	innerErr := ErrFileNotFound
	err := NewErrorWithWrapped(100, "wrapper", innerErr)

	if err.Code != 100 {
		t.Fatalf("Expected code 100, got %d", err.Code)
	}
	if err.Message != "wrapper" {
		t.Fatalf("Expected message 'wrapper', got '%s'", err.Message)
	}
	if err.Err != innerErr {
		t.Fatal("Wrapped error not preserved")
	}
}

func TestErrorString(t *testing.T) {
	err := ErrFileNotFound
	errStr := err.Error()

	if errStr == "" {
		t.Fatal("Error string should not be empty")
	}
	if errStr != "File not found" {
		t.Fatalf("Expected 'File not found', got '%s'", errStr)
	}
}

func TestErrorStringWithWrapped(t *testing.T) {
	innerErr := ErrFileOperation
	err := NewErrorWithWrapped(100, "wrapper message", innerErr)
	errStr := err.Error()

	if errStr == "" {
		t.Fatal("Error string should not be empty")
	}

	if len(errStr) < len("wrapper message") {
		t.Fatal("Wrapped error string should contain wrapper message")
	}
}

func TestErrorUnwrap(t *testing.T) {
	innerErr := ErrIndexOutOfRange
	err := NewErrorWithWrapped(100, "wrapper", innerErr)

	unwrapped := err.Unwrap()
	if unwrapped != innerErr {
		t.Fatal("Unwrap failed to return wrapped error")
	}
}

func TestErrorUnwrapNil(t *testing.T) {
	err := NewError(50, "no wrap")
	unwrapped := err.Unwrap()

	if unwrapped != nil {
		t.Fatal("Unwrap should return nil for non-wrapped error")
	}
}

func TestGetErrorCode(t *testing.T) {
	testCases := []struct {
		err      error
		expected int
	}{
		{ErrFileNotFound, ErrCodeFileNotFound},
		{ErrOutOfMemory, ErrCodeOutOfMemory},
		{ErrIndexOutOfRange, ErrCodeIndexOutOfRange},
		{ErrFileOperation, ErrCodeFileOperation},
	}

	for _, tc := range testCases {
		code := GetErrorCode(tc.err)
		if code != tc.expected {
			t.Fatalf("Expected code %d, got %d", tc.expected, code)
		}
	}
}

func TestGetErrorCodeUnknown(t *testing.T) {
	code := GetErrorCode(nil)
	if code != 0 {
		t.Fatalf("Expected 0 for nil error, got %d", code)
	}
}

func TestErrorCodes(t *testing.T) {
	if ErrCodeFileNotFound != -1 || ErrCodeOutOfMemory != -2 {
		t.Fatal("Error codes should have correct values")
	}
}
