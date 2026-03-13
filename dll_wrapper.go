package main

import (
	"C"
	"VirtualMemoryManagement/api"
	"unsafe"
)

func init() {
}

export func VMCreate(filename *C.char, size C.int, typ *C.char, stringLength C.int) C.int {
	filenameGo := C.GoString(filename)
	typGo := C.GoString(typ)
	result := api.VMCreate(filenameGo, int(size), typGo, int(stringLength))
	if result.IsSuccess() {
		return 1
	}
	return 0
}

export func VMOpen(filename *C.char) C.int {
	filenameGo := C.GoString(filename)
	result := api.VMOpen(filenameGo)
	if result.IsSuccess() {
		return 1
	}
	return 0
}

export func VMClose(handle C.int) C.int {
	result := api.VMClose(int(handle))
	if result.IsSuccess() {
		return 1
	}
	return 0
}

export func VMRead(handle C.int, index C.int, outBuffer *C.char, bufferSize C.int) C.int {
	result := api.VMRead(int(handle), int(index))
	if result.IsSuccess() {
		data := result.String()
		if len(data) >= int(bufferSize) {
			data = data[:int(bufferSize)-1]
		}
		C.memcpy(unsafe.Pointer(outBuffer), unsafe.Pointer(C.CString(data)), C.size_t(len(data)+1))
		return 1
	}
	return 0
}

export func VMWrite(handle C.int, index C.int, value *C.char) C.int {
	valueGo := C.GoString(value)
	result := api.VMWrite(int(handle), int(index), valueGo)
	if result.IsSuccess() {
		return 1
	}
	return 0
}

export func VMHelp(filename *C.char) {
	if filename != nil {
		filenameGo := C.GoString(filename)
		api.VMHelp(filenameGo)
	} else {
		api.VMHelp("")
	}
}

func main() {
}

