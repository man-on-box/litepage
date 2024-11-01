package build

import (
	"fmt"
	"time"

	"github.com/man-on-box/litepage/internal/common"
	"github.com/man-on-box/litepage/internal/file"
)

type SiteBuilder interface {
	Build() error
}

type siteBuilder struct {
	DistDir     string
	PublicDir   string
	Pages       *[]common.Page
	SiteDomain  string
	WithSitemap bool
}

func New(distDir string, publicDir string, pages *[]common.Page, siteDomain string, withSitemap bool) SiteBuilder {
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
	startTime := time.Now()

	err := file.CopyDir(b.PublicDir, b.DistDir)
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

	noOfPages := len(*b.Pages)
	pageStr := "page"
	if noOfPages > 1 {
		pageStr += "s"
	}

	fmt.Printf("Built %d %s in %.0f seconds\n", len(*b.Pages), pageStr, time.Since(startTime).Seconds())
	return nil
}

func (b *siteBuilder) createPages() error {
	for _, p := range *b.Pages {
		fmt.Printf("- creating %s...\n", p.Path)
		f, err := file.CreateFile(b.DistDir + p.Path)
		if err != nil {
			return err
		}
		p.Handler(f)
	}
	return nil
}

func (b *siteBuilder) createSitemap() error {
	f, err := file.CreateFile(b.DistDir + "/sitemap.xml")
	if err != nil {
		return err
	}
	sitemap := common.BuildSitemap(b.SiteDomain, b.Pages)
	_, err = f.Write([]byte(sitemap))
	return err
}
