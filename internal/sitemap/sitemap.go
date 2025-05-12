package sitemap

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/man-on-box/litepage/internal/model"
)

func BuildSitemap(domain string, pages *[]model.Page) string {
	var builder strings.Builder

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, p := range *pages {
		path := p.Path
		fileExt := filepath.Ext(path)
		isNotHTML := fileExt != ".html" && fileExt != ".htm"
		if isNotHTML {
			// do not include non html pages into sitemap
			continue
		}
		path = strings.TrimSuffix(path, filepath.Ext(path))
		path = strings.TrimSuffix(path, "index")
		builder.WriteString(fmt.Sprintf("<url><loc>https://%s%s</loc></url>", domain, path))
	}
	builder.WriteString("</urlset>")

	return builder.String()
}
