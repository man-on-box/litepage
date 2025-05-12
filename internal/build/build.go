package build

import (
	"fmt"
	"time"

	"github.com/man-on-box/litepage/internal/file"
	"github.com/man-on-box/litepage/internal/model"
	"github.com/man-on-box/litepage/internal/sitemap"
)

type SiteBuilder interface {
	Build() error
}

type Config struct {
	DistDir     string
	PublicDir   string
	Pages       *[]model.Page
	SiteDomain  string
	BasePath    string
	WithSitemap bool
}

type siteBuilder struct {
	Config Config
}

func New(config Config) SiteBuilder {
	b := &siteBuilder{
		Config: config,
	}
	return b
}

func (b *siteBuilder) Build() error {
	fmt.Printf("LITEPAGE building site '%s'...\n", b.Config.SiteDomain)
	startTime := time.Now()

	err := file.CopyDir(b.Config.PublicDir, b.Config.DistDir)
	if err != nil {
		return fmt.Errorf("Could not copy public directory: %w", err)
	}

	err = b.createPages()
	if err != nil {
		return fmt.Errorf("An error occurred while creating pages: %w", err)
	}

	if b.Config.WithSitemap {
		err = b.createSitemap()
		if err != nil {
			return fmt.Errorf("An error occurred while creating sitemap: %w", err)
		}
	}

	noOfPages := len(*b.Config.Pages)
	pageStr := "page"
	if noOfPages > 1 {
		pageStr += "s"
	}

	fmt.Printf("Built %d %s in %.0f seconds\n", len(*b.Config.Pages), pageStr, time.Since(startTime).Seconds())
	return nil
}

func (b *siteBuilder) createPages() error {
	for _, p := range *b.Config.Pages {
		fmt.Printf("- creating %s...\n", p.Path)
		f, err := file.CreateFile(b.Config.DistDir + p.Path)
		if err != nil {
			return err
		}
		p.Handler(f)
	}
	return nil
}

func (b *siteBuilder) createSitemap() error {
	f, err := file.CreateFile(b.Config.DistDir + "/sitemap.xml")
	if err != nil {
		return err
	}
	smap := sitemap.Build(b.Config.SiteDomain, b.Config.BasePath, b.Config.Pages)
	_, err = f.Write([]byte(smap))
	return err
}
