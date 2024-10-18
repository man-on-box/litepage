package common

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

func BuildSitemap(domain string, pageMap *PageMap) string {
	var builder strings.Builder
	sortedPaths := SortPageMapByPath(pageMap)

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, path := range sortedPaths {
		urlPath := fileToUrlPath(path)
		builder.WriteString(fmt.Sprintf("<url><loc>https://%s%s</loc></url>", domain, urlPath))
	}
	builder.WriteString("</urlset>")

	return builder.String()
}

func SortPageMapByPath(pageMap *PageMap) []string {
	keys := make([]string, len(*pageMap))
	i := 0
	for k := range *pageMap {
		keys[i] = k
		i++
	}
	slices.Sort(keys)
	return keys
}

func fileToUrlPath(path string) string {
	path = strings.TrimSuffix(path, filepath.Ext(path))
	path = strings.TrimSuffix(path, "index")
	return path
}
