package api

import (
	"VirtualMemoryManagement/config"
	"VirtualMemoryManagement/errors"
	"VirtualMemoryManagement/types/array"
	"VirtualMemoryManagement/types/result"
	"VirtualMemoryManagement/virtualmemory"
	"strconv"
	"sync"
)

var (
	mu        sync.Mutex
	handles   = make(map[int]*virtualmemory.VirtualArray)
	nextID    = 1
	cacheSize = config.DefaultCacheSize
)

func SetCacheSize(size int) {
	mu.Lock()
	defer mu.Unlock()

	if size < config.MinCacheSize {
		size = config.MinCacheSize
	}
	if size > config.MaxCacheSize {
		size = config.MaxCacheSize
	}
	cacheSize = size
}

func GetCacheSize() int {
	mu.Lock()
	defer mu.Unlock()
	return cacheSize
}

func VMCreate(filename string, size int, typ string, stringLength int) result.Result {
	mu.Lock()
	defer mu.Unlock()

	var arrayType array.Type
	switch typ {
	case "int", "I":
		arrayType = array.TypeInt
	case "char", "C":
		arrayType = array.TypeChar
	case "varchar", "V":
		arrayType = array.TypeVarchar
	default:
		return result.ErrorWithCode(errors.ErrCodeInvalidType, "Unknown type: "+typ)
	}

	va, err := virtualmemory.CreateWithCacheSize(filename, size, arrayType, stringLength, cacheSize)
	if err != nil {
		return result.Error(err)
	}

	id := nextID
	nextID++
	handles[id] = va

	return result.Success(filename)
}

func VMOpen(filename string) result.Result {
	mu.Lock()
	defer mu.Unlock()

	va, err := virtualmemory.OpenWithCacheSize(filename, cacheSize)
	if err != nil {
		return result.Error(err)
	}

	id := nextID
	nextID++
	handles[id] = va

	return result.Success(filename)
}

// VMClose ensures the handle is always removed from the map, even if errors occur
func VMClose(handle int) result.Result {
	mu.Lock()
	defer mu.Unlock()

	va, exists := handles[handle]
	if !exists {
		return result.ErrorWithCode(errors.ErrCodeInvalidHandle, "Invalid handle")
	}

	// Ensure handle is always removed, even if errors occur
	defer delete(handles, handle)

	if err := va.FlushDirtyPages(); err != nil {
		return result.Error(err)
	}

	if err := va.Close(); err != nil {
		return result.Error(err)
	}

	return result.Success("Closed")
}

func VMRead(handle int, index int) result.Result {
	mu.Lock()
	va, exists := handles[handle]
	mu.Unlock()

	if !exists {
		return result.ErrorWithCode(errors.ErrCodeInvalidHandle, "Invalid handle")
	}

	value, err := va.Read(index)
	if err != nil {
		return result.Error(err)
	}

	var resultStr string
	switch v := value.(type) {
	case int32:
		resultStr = strconv.FormatInt(int64(v), 10)
	case string:
		resultStr = v
	default:
		resultStr = ""
	}

	return result.Success(resultStr)
}

func VMWrite(handle int, index int, value string) result.Result {
	mu.Lock()
	va, exists := handles[handle]
	mu.Unlock()

	if !exists {
		return result.ErrorWithCode(errors.ErrCodeInvalidHandle, "Invalid handle")
	}

	arrayInfo := va.ArrayInfo()

	var writeValue interface{}
	switch arrayInfo.Type {
	case array.TypeInt:
		intVal, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			intVal = 0
		}
		writeValue = int32(intVal)
	case array.TypeChar, array.TypeVarchar:
		writeValue = value
	default:
		return result.ErrorWithCode(errors.ErrCodeInvalidType, "Invalid array type")
	}

	if err := va.Write(index, writeValue); err != nil {
		return result.Error(err)
	}

	return result.Success("Written")
}

func VMHelp(filename string) result.Result {
	help := `Virtual Memory Manager Commands:
Create <filename> <type> [<stringLength>] - Create new array
  Types: int, char(length), varchar(maxlength)
  Example: Create myfile.vm int
           Create myfile.vm char(50)
           Create myfile.vm varchar(100)

Open <filename> - Open existing array file

Close <handle> - Close array and flush to disk

Read <handle> <index> - Read element at index

Write <handle> <index> <value> - Write value to index
  String values must be quoted: "my string"

Stats <handle> - Show statistics

Help [<filename>] - Show this help

Exit - Close all and exit
`

	return result.Success(help)
}

func GetHandle() int {
	mu.Lock()
	defer mu.Unlock()

	for id := range handles {
		return id
	}
	return -1
}

func GetAllHandles() []int {
	mu.Lock()
	defer mu.Unlock()

	var ids []int
	for id := range handles {
		ids = append(ids, id)
	}
	return ids
}

func VMStats(handle int) result.Result {
	mu.Lock()
	va, exists := handles[handle]
	mu.Unlock()

	if !exists {
		return result.ErrorWithCode(errors.ErrCodeInvalidHandle, "Invalid handle")
	}

	stats := va.GetStats()
	return result.Success(stats)
}
