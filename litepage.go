package litepage

import (
	"errors"
	"os"

	"github.com/man-on-box/litepage/internal/build"
	"github.com/man-on-box/litepage/internal/serve"
	"github.com/man-on-box/litepage/pkg/types"
)

type Litepage interface {
	Run() error
	Serve(port string) error
	Page(filePath string, handler types.PageHandler)
}

type Option func(*litepage)

type litepage struct {
	siteDomain  string
	distDir     string
	publicDir   string
	withSitemap bool
	pages       *[]types.Page
}

func New(domain string, options ...Option) (Litepage, error) {
	lp := &litepage{
		siteDomain:  domain,
		distDir:     "dist",
		publicDir:   "public",
		withSitemap: true,
		pages:       &[]types.Page{},
	}

	for _, opt := range options {
		opt(lp)
	}

	if lp.siteDomain == "" {
		return nil, errors.New("site domain is required, please provide a domain like 'catpics.com'")
	}

	return lp, nil
}

func WithDistDir(distDis string) Option {
	return func(lp *litepage) {
		lp.distDir = distDis
	}
}

func WithPublicDir(publicDir string) Option {
	return func(lp *litepage) {
		lp.publicDir = publicDir
	}
}

func WithoutSitemap() Option {
	return func(lp *litepage) {
		lp.withSitemap = false
	}
}

func (lp *litepage) Page(filePath string, handler types.PageHandler) {
	*lp.pages = append(*lp.pages, types.Page{FilePath: filePath, Handler: handler})
}

// Run the litepage app. By default it will create the static site in the dist directory.
// If LP_MODE env variable is set to 'serve', it will instead serve the static site on port
// 3000, or on port specified if LP_PORT env variable was set.
func (lp *litepage) Run() error {
	mode := os.Getenv("LP_MODE")
	port := os.Getenv("LP_PORT")

	if mode == "serve" {
		return lp.Serve(port)
	} else {
		return lp.Build()
	}
}

func (lp *litepage) Build() error {
	builder := build.New(lp.distDir, lp.publicDir, lp.pages, lp.siteDomain, lp.withSitemap)
	return builder.Build()
}

func (lp *litepage) Serve(port string) error {
	server := serve.New(lp.publicDir, lp.pages)
	return server.Serve(port)
}
