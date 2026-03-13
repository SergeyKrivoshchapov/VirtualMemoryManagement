// Package storage provides interfaces and implementations for persistent storage operations.
package storage

import "VirtualMemoryManagement/types/page"

// PageStorage defines the interface for reading and writing pages to disk.
type PageStorage interface {
	ReadPage(pageNumber int) (*page.Page, error)
	WritePage(p *page.Page) error
	Close() error
	ArrayInfo() interface{}
}

// VarcharStorage defines the interface for reading and writing variable-length strings.
type VarcharStorage interface {
	WriteString(offset int64, s string) error
	ReadString(offset int64) (string, error)
	GetCurrentOffset() (int64, error)
	Close() error
}

