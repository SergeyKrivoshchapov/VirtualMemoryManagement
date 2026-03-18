package virtualmemory

import (
	"VirtualMemoryManagement/cache"
	"VirtualMemoryManagement/config"
	"VirtualMemoryManagement/errors"
	"VirtualMemoryManagement/storage"
	"VirtualMemoryManagement/types/array"
	"VirtualMemoryManagement/types/page"
	"encoding/binary"
	"strconv"
	"time"
)

type VirtualArray struct {
	pageStorage  storage.PageStorage
	varcharStore storage.VarcharStorage
	arrayInfo    *array.Info
	pageCache    cache.Cache
	varcharIndex map[int]int64
	cacheSize    int
}

// Create creates a new virtual array with the specified parameters.
// Uses DefaultCacheSize from config for cache capacity.
func Create(filename string, size int, typ array.Type, stringLength int) (*VirtualArray, error) {
	return CreateWithCacheSize(filename, size, typ, stringLength, config.DefaultCacheSize)
}

// CreateWithCacheSize creates a new virtual array with a custom cache size.
func CreateWithCacheSize(filename string, size int, typ array.Type, stringLength int, cacheSize int) (*VirtualArray, error) {
	if size <= 0 {
		return nil, errors.ErrInvalidType
	}

	if cacheSize < config.MinCacheSize {
		cacheSize = config.MinCacheSize
	}
	if cacheSize > config.MaxCacheSize {
		cacheSize = config.MaxCacheSize
	}

	pageFile := storage.NewPageFile(filename)
	if err := pageFile.Create(size, typ, stringLength); err != nil {
		return nil, err
	}

	va := &VirtualArray{
		pageStorage:  pageFile,
		arrayInfo:    pageFile.ArrayInfo(),
		pageCache:    cache.NewLRU(cacheSize),
		varcharIndex: make(map[int]int64),
		cacheSize:    cacheSize,
	}

	if typ == array.TypeVarchar {
		varcharFile := storage.NewVarcharFile(filename + ".varchar")
		if err := varcharFile.Create(); err != nil {
			return nil, errors.NewErrorWithWrapped(errors.ErrCodeFileOperation, "failed to create varchar file", err)
		}
		va.varcharStore = varcharFile
	}

	if err := va.loadInitialPages(); err != nil {
		return nil, err
	}

	return va, nil
}

// Open opens an existing virtual array from a file.
func Open(filename string) (*VirtualArray, error) {
	return OpenWithCacheSize(filename, config.DefaultCacheSize)
}

// OpenWithCacheSize opens an existing virtual array with a custom cache size.
func OpenWithCacheSize(filename string, cacheSize int) (*VirtualArray, error) {
	if cacheSize < config.MinCacheSize {
		cacheSize = config.MinCacheSize
	}
	if cacheSize > config.MaxCacheSize {
		cacheSize = config.MaxCacheSize
	}

	pageFile := storage.NewPageFile(filename)
	if err := pageFile.Open(filename); err != nil {
		return nil, err
	}

	va := &VirtualArray{
		pageStorage:  pageFile,
		arrayInfo:    pageFile.ArrayInfo(),
		pageCache:    cache.NewLRU(cacheSize),
		varcharIndex: make(map[int]int64),
		cacheSize:    cacheSize,
	}

	if va.arrayInfo.Type == array.TypeVarchar {
		varcharFile := storage.NewVarcharFile(filename + ".varchar")
		if err := varcharFile.Open(); err != nil {
			return nil, errors.NewErrorWithWrapped(errors.ErrCodeFileOperation, "failed to open varchar file", err)
		}
		va.varcharStore = varcharFile

		if err := va.loadVarcharIndex(); err != nil {
			return nil, err
		}
	}

	if err := va.loadInitialPages(); err != nil {
		return nil, err
	}

	return va, nil
}

// Close closes the virtual array and flushes all dirty pages.
func (va *VirtualArray) Close() error {
	if err := va.pageStorage.Close(); err != nil {
		return err
	}
	if va.varcharStore != nil {
		if err := va.varcharStore.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (va *VirtualArray) Read(index int) (interface{}, error) {
	if index < 0 || index >= va.arrayInfo.Size {
		return nil, errors.ErrIndexOutOfRange
	}

	p, offset, err := va.getPageAndOffset(index)
	if err != nil {
		return nil, err
	}

	switch va.arrayInfo.Type {
	case array.TypeInt:
		// Check if element was ever written (bit set in bitmap)
		bitIndex := index % config.BitsPerPage
		isBitSet, err := p.IsBitSet(bitIndex)
		if err != nil {
			return nil, err
		}
		if !isBitSet {
			return nil, errors.NewError(errors.ErrCodeIndexOutOfRange, "Element at index "+strconv.Itoa(index)+" has not been written")
		}
		return va.readInt(p.Data(), offset), nil

	case array.TypeChar:
		// Check if element was ever written (bit set in bitmap)
		bitIndex := index % config.BitsPerPage
		isBitSet, err := p.IsBitSet(bitIndex)
		if err != nil {
			return nil, err
		}
		if !isBitSet {
			return nil, errors.NewError(errors.ErrCodeIndexOutOfRange, "Element at index "+strconv.Itoa(index)+" has not been written")
		}
		return va.readChar(p.Data(), offset, va.arrayInfo.StringLength), nil

	case array.TypeVarchar:
		// For varchar, offset=0 means "not written"
		varcharOffset := int64(va.readInt(p.Data(), offset))
		if varcharOffset == 0 {
			return nil, errors.NewError(errors.ErrCodeIndexOutOfRange, "Element at index "+strconv.Itoa(index)+" has not been written")
		}
		if va.varcharStore == nil {
			return "", errors.ErrFileOperation
		}
		str, err := va.varcharStore.ReadString(varcharOffset)
		if err != nil {
			return "", nil
		}
		return str, nil

	default:
		return nil, errors.ErrInvalidType
	}
}

func (va *VirtualArray) Write(index int, value interface{}) error {
	if index < 0 || index >= va.arrayInfo.Size {
		return errors.ErrIndexOutOfRange
	}

	p, offset, err := va.getPageAndOffset(index)
	if err != nil {
		return err
	}

	switch va.arrayInfo.Type {
	case array.TypeInt:
		intVal, ok := value.(int32)
		if !ok {
			return errors.ErrInvalidType
		}
		va.writeInt(p.Data(), offset, intVal)

	case array.TypeChar:
		strVal, ok := value.(string)
		if !ok {
			return errors.ErrInvalidType
		}
		va.writeChar(p.Data(), offset, strVal, va.arrayInfo.StringLength)

	case array.TypeVarchar:
		strVal, ok := value.(string)
		if !ok {
			return errors.ErrInvalidType
		}

		varcharOffset, err := va.varcharStore.GetCurrentOffset()
		if err != nil {
			return err
		}

		if err := va.varcharStore.WriteString(varcharOffset, strVal); err != nil {
			return err
		}

		va.writeInt(p.Data(), offset, int32(varcharOffset))
		va.varcharIndex[index] = varcharOffset

	default:
		return errors.ErrInvalidType
	}

	bitIndex := index % config.BitsPerPage

	if err := p.SetBit(bitIndex); err != nil {
		return err
	}
	p.MakeDirty()

	return nil
}

func (va *VirtualArray) ArrayInfo() *array.Info {
	return va.arrayInfo
}

func (va *VirtualArray) getPageAndOffset(index int) (*page.Page, int, error) {
	pageNumber := index / config.BitsPerPage
	offsetInPage := (index % config.BitsPerPage) * va.arrayInfo.ElementSize

	p, err := va.ensurePageInCache(pageNumber)
	if err != nil {
		return nil, 0, err
	}

	return p, offsetInPage, nil
}

func (va *VirtualArray) ensurePageInCache(pageNumber int) (*page.Page, error) {
	if p := va.pageCache.Get(pageNumber); p != nil {
		p.MarkAccessed()
		return p, nil
	}

	p, err := va.pageStorage.ReadPage(pageNumber)
	if err != nil {
		return nil, err
	}

	evicted := va.pageCache.Put(p)
	if evicted != nil && evicted.Dirty {
		if err := va.pageStorage.WritePage(evicted); err != nil {
			return nil, err
		}
	}

	p.MarkAccessed()
	return p, nil
}

func (va *VirtualArray) loadInitialPages() error {
	count := va.cacheSize
	if count > va.arrayInfo.PageCount {
		count = va.arrayInfo.PageCount
	}

	for i := 0; i < count; i++ {
		if p := va.pageCache.Get(i); p != nil {
			continue
		}

		p, err := va.pageStorage.ReadPage(i)
		if err != nil {
			// Skip pages that can't be read - they'll be loaded on demand
			// or newly created pages if they're accessed
			continue
		}

		va.pageCache.Put(p)
	}

	return nil
}

func (va *VirtualArray) loadVarcharIndex() error {
	for i := 0; i < va.arrayInfo.Size; i++ {
		if p := va.pageCache.Get(i / config.BitsPerPage); p != nil {
			offset := (i % config.BitsPerPage) * va.arrayInfo.ElementSize
			varcharOffset := int64(va.readInt(p.Data(), offset))
			if varcharOffset != 0 {
				va.varcharIndex[i] = varcharOffset
			}
		}
	}
	return nil
}

func (va *VirtualArray) readInt(data []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
}

func (va *VirtualArray) writeInt(data []byte, offset int, value int32) {
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(value))
}

func (va *VirtualArray) readInt64(data []byte, offset int) int64 {
	return int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
}

func (va *VirtualArray) writeInt64(data []byte, offset int, value int64) {
	binary.LittleEndian.PutUint64(data[offset:offset+8], uint64(value))
}

func (va *VirtualArray) readChar(data []byte, offset int, length int) string {
	endIndex := offset + length
	if endIndex > len(data) {
		endIndex = len(data)
	}

	str := string(data[offset:endIndex])
	if idx := -1; true {
		for i, b := range data[offset:endIndex] {
			if b == 0 {
				idx = i
				break
			}
		}
		if idx != -1 {
			str = string(data[offset : offset+idx])
		}
	}
	return str
}

func (va *VirtualArray) writeChar(data []byte, offset int, value string, length int) {
	for i := 0; i < length && i < len(value); i++ {
		data[offset+i] = value[i]
	}
	for i := len(value); i < length; i++ {
		data[offset+i] = 0
	}
}

// FlushDirtyPages writes all dirty pages back to storage.
func (va *VirtualArray) FlushDirtyPages() error {
	pages := va.pageCache.All()
	for _, p := range pages {
		if p.Dirty {
			if err := va.pageStorage.WritePage(p); err != nil {
				return err
			}
			p.Dirty = false
		}
	}

	// Additional sync to ensure all data is written to disk
	if err := va.pageStorage.Sync(); err != nil {
		return err
	}

	return nil
}

func (va *VirtualArray) GetStats() string {
	pages := va.pageCache.All()
	stats := "Virtual Array Stats:\n"
	stats += "Array Type: " + va.arrayInfo.Type.String() + "\n"
	stats += "Array Size: " + strconv.Itoa(va.arrayInfo.Size) + "\n"
	stats += "Element Size: " + strconv.Itoa(va.arrayInfo.ElementSize) + " bytes\n"
	stats += "Total Pages: " + strconv.Itoa(va.arrayInfo.PageCount) + "\n"
	stats += "Cached Pages: " + strconv.Itoa(len(pages)) + "\n"
	stats += "Page Details:\n"

	for _, p := range pages {
		stats += "  Page " + strconv.Itoa(p.AbsoluteNumber) +
			" | Dirty: " + strconv.FormatBool(p.Dirty) +
			" | Accesses: " + strconv.Itoa(p.AccessCounter) +
			" | Last Access: " + p.AccessTime.Format(time.RFC3339Nano) + "\n"
	}

	return stats
}
