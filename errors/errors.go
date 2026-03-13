package errors

import (
	"errors"
	"fmt"
)

type VMMError struct {
	Code    int
	Message string
	Err     error
}

func (e *VMMError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap enables errors.Is and errors.As
func (e *VMMError) Unwrap() error {
	return e.Err
}

const (
	ErrCodeFileNotFound     = -1
	ErrCodeOutOfMemory      = -2
	ErrCodeIndexOutOfRange  = -3
	ErrCodeFileOperation    = -4
	ErrCodeInvalidType      = -5
	ErrCodeInsufficientDisk = -6
	ErrCodeInvalidHandle    = -7
	ErrCodePageNotFound     = -8
)

var (
	ErrFileNotFound     = &VMMError{Code: ErrCodeFileNotFound, Message: "File not found"}
	ErrOutOfMemory      = &VMMError{Code: ErrCodeOutOfMemory, Message: "Out of memory"}
	ErrIndexOutOfRange  = &VMMError{Code: ErrCodeIndexOutOfRange, Message: "Index out of range"}
	ErrFileOperation    = &VMMError{Code: ErrCodeFileOperation, Message: "File operation failed"}
	ErrInvalidType      = &VMMError{Code: ErrCodeInvalidType, Message: "Invalid array type"}
	ErrInsufficientDisk = &VMMError{Code: ErrCodeInsufficientDisk, Message: "Insufficient disk space"}
	ErrInvalidHandle    = &VMMError{Code: ErrCodeInvalidHandle, Message: "Invalid handle"}
	ErrPageNotFound     = &VMMError{Code: ErrCodePageNotFound, Message: "Page not found"}
)

// NewError creates a custom error
func NewError(code int, message string) *VMMError {
	return &VMMError{Code: code, Message: message}
}

func NewErrorWithWrapped(code int, message string, err error) *VMMError {
	return &VMMError{Code: code, Message: message, Err: err}
}

func GetErrorCode(err error) int {
	if err == nil {
		return 0
	}
	var vmmErr *VMMError
	if errors.As(err, &vmmErr) {
		return vmmErr.Code
	}
	return -999
}
