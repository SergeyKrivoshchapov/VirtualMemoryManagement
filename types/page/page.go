package page

import (
	"VirtualMemoryManagement/config"
	"VirtualMemoryManagement/types/bitmap"
	"time"
)

type Page struct {
	AbsoluteNumber int
	Dirty          bool
	AccessTime     time.Time
	AccessCounter  int
	WriteProtected bool
	bitmap         *bitmap.BitMap
	data           []byte
}

// New creates a new page without data
func New(number int, elemSize int) *Page {
	dataSize := config.PageDataSize(elemSize)
	return &Page{
		AbsoluteNumber: number,
		Dirty:          false,
		AccessTime:     time.Now(),
		bitmap:         bitmap.New(),
		data:           make([]byte, dataSize),
	}
}

// NewWithData creates a page and loads data from byte slices
func NewWithData(number int, elemSize int, bitmapData []byte, pageData []byte) (*Page, error) {
	p := New(number, elemSize)
	if err := p.bitmap.FromBytes(bitmapData); err != nil {
		return nil, err
	}
	copy(p.data, pageData)
	return p, nil
}

func (p *Page) Bitmap() *bitmap.BitMap {
	return p.bitmap
}

func (p *Page) Data() []byte {
	return p.data
}

func (p *Page) MarkAccessed() {
	p.AccessTime = time.Now()
	p.AccessCounter++
}

func (p *Page) MakeDirty() {
	p.Dirty = true
	p.MarkAccessed()
}

func (p *Page) SetBit(pos int) error {
	if err := p.bitmap.Set(pos); err != nil {
		return err
	}
	p.Dirty = true
	return nil
}

func (p *Page) IsBitSet(pos int) (bool, error) {
	return p.bitmap.IsSet(pos)
}

func (p *Page) DataSize() int {
	return len(p.data)
}
