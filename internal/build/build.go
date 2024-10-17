package build

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/man-on-box/litepage/internal/file"
	"github.com/man-on-box/litepage/pkg/types"
)

type SiteBuilder interface {
	Build() error
}

type siteBuilder struct {
	DistDir     string
	PublicDir   string
	Pages       *[]types.Page
	SiteDomain  string
	WithSitemap bool
}

func New(distDir string, publicDir string, pages *[]types.Page, siteDomain string, withSitemap bool) SiteBuilder {
	b := &siteBuilder{
		DistDir:     distDir,
		PublicDir:   publicDir,
		Pages:       pages,
		SiteDomain:  siteDomain,
		WithSitemap: withSitemap,
	}
	return b
}

func (b *siteBuilder) Build() error {
	fmt.Printf("LITEPAGE building site '%s'...\n", b.SiteDomain)
	err := file.RemoveDir(b.DistDir)
	if err != nil {
		return fmt.Errorf("Could not remove dist directory: %w", err)
	}

	err = file.CopyDir(b.PublicDir, b.DistDir)
	if err != nil {
		return fmt.Errorf("Could not copy public directory: %w", err)
	}

	err = b.createPages()
	if err != nil {
		return fmt.Errorf("An error occurred while creating pages: %w", err)
	}

	if b.WithSitemap {
		err = b.createSitemap()
		if err != nil {
			return fmt.Errorf("An error occurred while creating sitemap: %w", err)
		}
	}

	return nil
}

func (b *siteBuilder) createPages() error {
	for _, page := range *b.Pages {
		fmt.Printf("- creating %s...\n", page.FilePath)
		f, err := file.CreateFile(b.DistDir + page.FilePath)
		if err != nil {
			return err
		}
		page.Handler(f)
	}
	return nil
}

func (b *siteBuilder) createSitemap() error {
	f, err := file.CreateFile(b.DistDir + "/sitemap.xml")
	if err != nil {
		return err
	}
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, page := range *b.Pages {
		urlPath := b.fileToUrlPath(page.FilePath)
		builder.WriteString(fmt.Sprintf("<url><loc>https://%s%s</loc></url>", b.SiteDomain, urlPath))
	}
	builder.WriteString("</urlset>")
	_, err = f.Write([]byte(builder.String()))
	return err
}

func (b *siteBuilder) fileToUrlPath(path string) string {
	path = strings.TrimSuffix(path, filepath.Ext(path))
	path = strings.TrimSuffix(path, "index")
	return path
}
