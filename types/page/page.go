package page

import (
	"VirtualMemoryManagement/config"
	"VirtualMemoryManagement/types/bitmap"
	"time"
)

// Page страница в памяти
type Page struct {
	AbsoluteNumber int       // абс номер страницы в файле
	Dirty          bool      // флаг модификации
	AccessTime     time.Time // последнее обращение
	AccessCounter  int       // кол-во обращений
	WriteProtected bool      // защита от записи
	bitmap         *bitmap.BitMap
	data           []byte // данные страницы
}

// New Создать без данных
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

// NewWithData Создать с готовыми данными
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
