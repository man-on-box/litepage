package build

import (
	"fmt"

	"github.com/man-on-box/litepage/internal/common"
	"github.com/man-on-box/litepage/internal/file"
)

type SiteBuilder interface {
	Build() error
}

type siteBuilder struct {
	DistDir     string
	PublicDir   string
	PageMap     *common.PageMap
	SiteDomain  string
	WithSitemap bool
}

func New(distDir string, publicDir string, pageMap *common.PageMap, siteDomain string, withSitemap bool) SiteBuilder {
	b := &siteBuilder{
		DistDir:     distDir,
		PublicDir:   publicDir,
		PageMap:     pageMap,
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
	for path, handler := range *b.PageMap {
		fmt.Printf("- creating %s...\n", path)
		f, err := file.CreateFile(b.DistDir + path)
		if err != nil {
			return err
		}
		handler(f)
	}
	return nil
}

func (b *siteBuilder) createSitemap() error {
	f, err := file.CreateFile(b.DistDir + "/sitemap.xml")
	if err != nil {
		return err
	}
	sitemap := common.BuildSitemap(b.SiteDomain, b.PageMap)
	_, err = f.Write([]byte(sitemap))
	return err
}
