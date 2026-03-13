package array

import "VirtualMemoryManagement/config"

type Type byte

const (
	TypeInt     Type = 'I' // Целые числа
	TypeChar    Type = 'C' // Строки фиксированной длины
	TypeVarchar Type = 'V' // Строки переменной длины
)

func (t Type) String() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeChar:
		return "char"
	case TypeVarchar:
		return "varchar"
	default:
		return "unknown"
	}
}

type Info struct {
	Size         int // количество элементов массива
	Type         Type
	StringLength int // для char и varchar
	elementSize  int
	pageCount    int
}

func NewInfo(size int, typ Type, stringLength int) *Info {
	elemSize := calcElementSize(typ, stringLength)

	return &Info{
		Size:         size,
		Type:         typ,
		StringLength: stringLength,
		ElementSize:  elemSize,
		PageCount:    (size + config.BitsPerPage - 1) / config.BitsPerPage,
	}
}

func calcElementSize(typ Type, stringLength int) int {
	switch typ {
	case TypeInt:
		return 4
	case TypeChar:
		return stringLength
	case TypeVarchar:
		return 4 // указатель
	default:
		return 0
	}
}
