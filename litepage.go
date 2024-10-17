// Package litepage is a minimalist, zero dependency static
// site builder for Go, to help move towards a simpler web.
//
// Create your pages with the Page() method to build out your site.
// Put your static assets like images, JS or CSS in the /public
// directory.
//
// You can then use Build() to build the static site into your
// dist directory, or Serve() to host the site locally, useful
// for local development.
//
// For convenience, you can use BuildOrServe() method to either
// build or serve, depending on environment variables. This allows
// you to serve locally, while build for production.
package litepage

import (
	"errors"
	"os"

	"github.com/man-on-box/litepage/internal/build"
	"github.com/man-on-box/litepage/internal/serve"
	"github.com/man-on-box/litepage/pkg/types"
)

type Litepage interface {
	// Build generates the static site in the dist directory using assets
	// from the public directory and pages registered in litepage instance,
	// ready to be hosted on your static site hosting service.
	Build() error
	// Serve hosts the static site on the specified port, instead of
	// writing the site to the dist directory. Use this for local development.
	Serve(port string) error
	// BuildOrServe by default will build the static site in the dist directory.
	// If LP_MODE env variable is set to 'serve', it will instead serve the static site on port
	// 3000, or on port specified if LP_PORT env variable was set.
	//
	// This is useful for when you want to serve the site during development and
	// build it for production without requiring changes to the code.
	BuildOrServe() error
	// Page registers a new page with the specified relative file path and handler.
	// Create a page by specifying the relative path the page should be created,
	// as well as the handler to render the page contents to the standard writer
	// interface.
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

// New creates a new Litepage instance with the specified domain and optional configurations.
// The domain parameter is required and must be a non-empty string representing the site's domain.
// The options parameter allows for additional configurations to be applied to the Litepage instance.
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

func (lp *litepage) Serve(port string) error {
	server := serve.New(lp.publicDir, lp.pages, lp.siteDomain, lp.withSitemap)
	return server.Serve(port)
}

func (lp *litepage) Build() error {
	builder := build.New(lp.distDir, lp.publicDir, lp.pages, lp.siteDomain, lp.withSitemap)
	return builder.Build()
}

func (lp *litepage) BuildOrServe() error {
	mode := os.Getenv("LP_MODE")
	port := os.Getenv("LP_PORT")

	if mode == "serve" {
		return lp.Serve(port)
	} else {
		return lp.Build()
	}
}
