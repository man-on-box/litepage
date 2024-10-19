package common

import (
	"fmt"
	"path/filepath"
	"strings"
)

func BuildSitemap(domain string, pages *[]Page) string {
	var builder strings.Builder

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, p := range *pages {
		urlPath := fileToUrlPath(p.Path)
		builder.WriteString(fmt.Sprintf("<url><loc>https://%s%s</loc></url>", domain, urlPath))
	}
	builder.WriteString("</urlset>")

	return builder.String()
}

func fileToUrlPath(path string) string {
	path = strings.TrimSuffix(path, filepath.Ext(path))
	path = strings.TrimSuffix(path, "index")
	return path
}
