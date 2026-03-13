// +build dll

package main

import (
	"C"
	"VirtualMemoryManagement/api"
)

//export VMCreate
func VMCreate(filename *C.char, size C.int, typ *C.char, stringLength C.int) C.int {
	filenameGo := C.GoString(filename)
	typGo := C.GoString(typ)
	result := api.VMCreate(filenameGo, int(size), typGo, int(stringLength))
	if result.IsSuccess() {
		return 1
	}
	return 0
}

//export VMOpen
func VMOpen(filename *C.char) C.int {
	filenameGo := C.GoString(filename)
	result := api.VMOpen(filenameGo)
	if result.IsSuccess() {
		return 1
	}
	return 0
}

//export VMClose
func VMClose(handle C.int) C.int {
	result := api.VMClose(int(handle))
	if result.IsSuccess() {
		return 1
	}
	return 0
}

//export VMRead
func VMRead(handle C.int, index C.int) C.Result {
	result := api.VMRead(int(handle), int(index))
	cResult := C.Result{}
	cResult.success = C.int32_t(result.Success)
	cResult.error_code = C.int32_t(result.ErrorCode)
	for i := 0; i < len(result.Data) && i < 256; i++ {
		cResult.data[i] = C.char(result.Data[i])
	}
	return cResult
}

//export VMWrite
func VMWrite(handle C.int, index C.int, value *C.char) C.int {
	valueGo := C.GoString(value)
	result := api.VMWrite(int(handle), int(index), valueGo)
	if result.IsSuccess() {
		return 1
	}
	return 0
}

//export VMHelp
func VMHelp(filename *C.char) C.Result {
	var filenameGo string
	if filename != nil {
		filenameGo = C.GoString(filename)
	}
	result := api.VMHelp(filenameGo)
	cResult := C.Result{}
	cResult.success = C.int32_t(result.Success)
	cResult.error_code = C.int32_t(result.ErrorCode)
	for i := 0; i < len(result.Data) && i < 256; i++ {
		cResult.data[i] = C.char(result.Data[i])
	}
	return cResult
}




