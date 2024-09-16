package pages

import (
	"io"
)

func (p *Pages) setupPages() *[]Page {
	pages := []Page{
		{
			filePath: "/index.html",
			render: func(w io.Writer) {
				d := p.data.NewPageIndex()
				p.executeTemplate(w, "page-index", d)
			},
		},
	}

	return &pages
}
