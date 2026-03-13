package result

import (
	"VirtualMemoryManagement/errors"
)

// Result единый формат возврата для API
type Result struct {
	Success   int // 1 - успех, 0 - ошибка
	Data      [256]byte
	ErrorCode int // код ошибки
}

func Success(data string) Result {
	var r Result
	r.Success = 1
	copy(r.Data[:], data)
	return r
}

func Error(err error) Result {
	var r Result
	r.Success = 0
	code := errors.GetErrorCode(err)
	r.ErrorCode = int(code)
	msg := err.Error()
	copy(r.Data[:], msg)

	return r
}

func ErrorWithCode(code int, message string) Result {
	var r Result
	r.Success = 0
	r.ErrorCode = code
	copy(r.Data[:], message)
	return r
}

func (r *Result) String() string {
	for i, b := range r.Data {
		if b == 0 {
			return string(r.Data[:i])
		}
	}
	return string(r.Data[:])
}

func (r *Result) IsSuccess() bool {
	return r.Success == 1
}

func (r *Result) GetErrorMessage() string {
	if r.IsSuccess() {
		return ""
	}
	return r.String()
}
