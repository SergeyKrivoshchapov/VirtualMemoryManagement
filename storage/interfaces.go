package storage

import (
	"VirtualMemoryManagement/types/array"
	"VirtualMemoryManagement/types/page"
)

type PageStorage interface {
	ReadPage(pageNumber int) (*page.Page, error)
	WritePage(p *page.Page) error
	Close() error
	ArrayInfo() *array.Info
}

type VarcharStorage interface {
	WriteString(offset int64, s string) error
	ReadString(offset int64) (string, error)
	GetCurrentOffset() (int64, error)
	Close() error
}
