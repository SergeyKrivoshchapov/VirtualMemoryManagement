//go:build dll
// +build dll

package main

/*
#include <stdint.h>

typedef struct {
    int32_t success;
    char data[256];
    int32_t error_code;
} Result;
*/
import "C"
import (
	"VirtualMemoryManagement/api"
	"VirtualMemoryManagement/types/result"
)

func toGoResultC(r result.Result) C.Result {
	cResult := C.Result{}
	cResult.success = C.int32_t(r.Success)
	cResult.error_code = C.int32_t(r.ErrorCode)
	for i := 0; i < len(r.Data) && i < 256; i++ {
		cResult.data[i] = C.char(r.Data[i])
	}
	return cResult
}

//export VMCreate
func VMCreate(filename *C.char, size C.int, typ *C.char, stringLength C.int) C.Result {
	filenameGo := C.GoString(filename)
	typGo := C.GoString(typ)
	return toGoResultC(api.VMCreate(filenameGo, int(size), typGo, int(stringLength)))
}

//export VMOpen
func VMOpen(filename *C.char) C.Result {
	filenameGo := C.GoString(filename)
	return toGoResultC(api.VMOpen(filenameGo))
}

//export VMRead
func VMRead(handle C.int, index C.int) C.Result {
	return toGoResultC(api.VMRead(int(handle), int(index)))
}

//export VMWrite
func VMWrite(handle C.int, index C.int, value *C.char) C.Result {
	valueGo := C.GoString(value)
	return toGoResultC(api.VMWrite(int(handle), int(index), valueGo))
}

//export VMHelp
func VMHelp(filename *C.char, helpText *C.char) C.Result {
	var filenameGo string
	if filename != nil {
		filenameGo = C.GoString(filename)
	}
	var helpTextGo string
	if helpText != nil {
		helpTextGo = C.GoString(helpText)
	}
	return toGoResultC(api.VMHelp(filenameGo, helpTextGo))
}

//export VMStats
func VMStats(handle C.int) C.Result {
	return toGoResultC(api.VMStats(int(handle)))
}
