package config

const (
	BitsPerPage      = 128             // 128 элементов на страницу
	BytesPerBitmap   = BitsPerPage / 8 // bitmap в байтах
	PhysicalPageSize = 512             // 512 байт - физ размер страницы по ТЗ
)

// PageDataSize размер данных страницы для типа
func PageDataSize(elemSize int) int {
	return BitsPerPage * elemSize
}

func TotalPageSize(elemSize int) int {
	dataSize := PageDataSize(elemSize)
	total := BytesPerBitmap + dataSize
	return ((total + PhysicalPageSize - 1) / PhysicalPageSize) * PhysicalPageSize // Выравниваем до PhysicalPageSize
}
