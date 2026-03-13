package array

type Type byte

const (
	TypeInt     Type = 'I' // Целые числа
	TypeChar    Type = 'C' // Строки фиксированной длины
	TypeVarchar Type = 'V' // Строки переменной длины
)

func (t ArrayType) String() string {
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

type ArrayInfo struct {
	Size            int // количество элементов массива
	Type            ArrayType
	StringLength    int // для char и varchar
	ElementSize     int
	PageCount       int
	PageSize        int
	ElementsPerPage int
}
