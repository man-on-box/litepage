package litepage

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/man-on-box/litepage/internal/util"
)

func (lp *LitePage) Build() error {
	fmt.Printf("LITEPAGE building site '%s'...\n", lp.config.siteDomain)

	err := util.RemoveDir(lp.config.distDir)
	if err != nil {
		return fmt.Errorf("Could not remove dist directory: %w", err)
	}

	err = util.CopyDir(lp.config.publicDir, lp.config.distDir)
	if err != nil {
		return fmt.Errorf("Could not copy public directory: %w", err)
	}

	err = lp.createPages()
	if err != nil {
		return fmt.Errorf("An error occurred while creating pages: %w", err)
	}

	if lp.config.withSitemap {
		err = lp.createSitemap()
		if err != nil {
			return fmt.Errorf("An error occurred while creating sitemap: %w", err)
		}
	}

	return nil
}

func (lp *LitePage) createPages() error {
	for _, page := range *lp.pages {
		fmt.Printf("- creating %s...\n", page.filePath)
		f, err := util.CreateFile(lp.config.distDir + page.filePath)
		if err != nil {
			return err
		}
		page.handler(f)
	}
	return nil
}

func (lp *LitePage) createSitemap() error {
	f, err := util.CreateFile(lp.config.distDir + "/sitemap.xml")
	if err != nil {
		return err
	}
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, page := range *lp.pages {
		urlPath := fileToUrlPath(page.filePath)
		builder.WriteString(fmt.Sprintf(`
        <url>
		<loc>https://%s%s</loc>
        </url>`, lp.config.siteDomain, urlPath))
	}
	builder.WriteString(`</urlset>`)
	_, err = f.Write([]byte(builder.String()))
	return err
}

func fileToUrlPath(path string) string {
	path = strings.TrimSuffix(path, filepath.Ext(path))
	path = strings.TrimSuffix(path, "index")
	return path
}
