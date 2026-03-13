package errors

import "errors"

type VMMError struct {
	Code    int
	Message string
}

func (e *VMMError) Error() string {
	return e.Message
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
	ErrOutOfMemory      = &VMMError{Code: ErrCodeOutOfMemory, Message: "Out of use"}
	ErrIndexOutOfRange  = &VMMError{Code: ErrCodeIndexOutOfRange, Message: "Index out of range"}
	ErrFileOperation    = &VMMError{Code: ErrCodeFileOperation, Message: "File operation failed"}
	ErrInvalidType      = &VMMError{Code: ErrCodeInvalidType, Message: "Invalid array type"}
	ErrInsufficientDisk = &VMMError{Code: ErrCodeInsufficientDisk, Message: "Insufficient disk space"}
	ErrInvalidHandle    = &VMMError{Code: ErrCodeInvalidHandle, Message: "Invalid handle"}
	ErrPageNotFound     = &VMMError{Code: ErrCodePageNotFound, Message: "Page not found"}
)

// NewError Ошибка произвольного содержимого
func NewError(code int, message string) *VMMError {
	return &VMMError{Code: code, Message: message}
}

// GetErrorCode Вернет код ошибки
func GetErrorCode(err error) int {
	if err == nil {
		return 0
	}
	if vmmErr, ok := errors.AsType[*VMMError](err); ok {
		return vmmErr.Code
	}
	return -999
}
