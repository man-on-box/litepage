package common

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/man-on-box/litepage/pkg/types"
)

func BuildSitemap(domain string, pages *[]types.Page) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, page := range *pages {
		urlPath := fileToUrlPath(page.FilePath)
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
