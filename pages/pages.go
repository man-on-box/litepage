package pages

import (
	"fmt"
	"io"
)

func (p *Pages) setupPages() *[]Page {
	fmt.Println("Setting up pages")
	pages := []Page{
		{
			filePath: "/index.html",
			render: func(w io.Writer) {
				d := p.data.NewPageIndex()
				p.executeTemplate(w, "page-index", d)
			},
		},
		{
			filePath: "/test.html",
			render: func(w io.Writer) {
				d := p.data.NewPageIndex()
				p.executeTemplate(w, "page-index", d)
			},
		},
	}

	return &pages
}
