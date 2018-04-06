package core

import (
	"strings"
)

// Paginator paginates lines of text into bigger messages.
type Paginator struct {
	MaxSize int

	pages []string
	page *strings.Builder
}

// NewPaginator creates a new Paginator.
func NewPaginator(maxSize uint) Paginator {
	return Paginator{
		MaxSize: int(maxSize),
		pages: make([]string, 0, 10),
		page: &strings.Builder{},
	}
}

// AddLine adds a line to be paginated.
func (p *Paginator) AddLine(line string) {
	if p.page.Len() + len(line) + 1 > p.MaxSize {
		p.EndPage()
	}

	p.page.WriteString(line)
	p.page.WriteByte('\n')
}

// EndPage ends the current page.
func (p *Paginator) EndPage() {
	if p.page.Len() != 0 {
		p.pages = append(p.pages, p.page.String())
		p.page.Reset()
	}
}

// Pages returns the finalized pages. May be called multiple times.
func (p *Paginator) Pages() []string {
	p.EndPage()
	return p.pages
}
