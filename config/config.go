package config

const (
	BitsPerPage      = 128
	BytesPerBitmap   = BitsPerPage / 8
	PhysicalPageSize = 512

	MinCacheSize     = 3
	MaxCacheSize     = 100
	DefaultCacheSize = 10
)

func PageDataSize(elemSize int) int {
	return BitsPerPage * elemSize
}

func TotalPageSize(elemSize int) int {
	dataSize := PageDataSize(elemSize)
	total := BytesPerBitmap + dataSize
	return ((total + PhysicalPageSize - 1) / PhysicalPageSize) * PhysicalPageSize
}
