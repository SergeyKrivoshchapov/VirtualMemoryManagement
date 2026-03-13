package page

import (
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
	Data           []byte // данные страницы
}

func NewPage(number int, dataSize int) *Page {
	return &Page{
		AbsoluteNumber: number,
		Dirty:          false,
		AccessTime:     time.Now(),
		AccessCounter:  0,
		WriteProtected: false,
		bitmap:         bitmap.New(),
		Data:           make([]byte, dataSize),
	}
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

func (p *Page) BitmapBytes() []byte {
	return p.bitmap.Bytes()
}
