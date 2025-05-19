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
	"fmt"
	"io"
	"os"

	"github.com/man-on-box/litepage/internal/build"
	"github.com/man-on-box/litepage/internal/model"
	"github.com/man-on-box/litepage/internal/serve"
	"github.com/man-on-box/litepage/internal/validate"
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
	Page(filePath string, handler func(w io.Writer)) error
}

type Option func(*litepage) error

type litepage struct {
	siteDomain  string
	distDir     string
	publicDir   string
	basePath    string
	withSitemap bool
	pages       *[]model.Page
	pathMap     map[string]bool
}

// New creates a new Litepage instance with the specified domain and optional configurations.
// The domain parameter is required and must be a non-empty string representing the site's domain.
// The options parameter allows for additional configurations to be applied to the Litepage instance.
func New(domain string, options ...Option) (Litepage, error) {
	if domain == "" {
		return nil, fmt.Errorf("site domain is required, please provide a domain like 'catpics.com'")
	}
	err := validate.IsValidDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("site domain is not valid: %v", err)
	}

	lp := &litepage{
		siteDomain:  domain,
		distDir:     "dist",
		publicDir:   "public",
		withSitemap: true,
		pathMap:     map[string]bool{},
		pages:       &[]model.Page{},
	}

	for _, opt := range options {
		if err := opt(lp); err != nil {
			return nil, err
		}
	}

	return lp, nil
}

// Specify the directory to place your site when building. By default it is "dist".
func WithDistDir(distDis string) Option {
	return func(lp *litepage) error {
		lp.distDir = distDis
		return nil
	}
}

// Specify the directory where your public assets are. By default it is "public".
func WithPublicDir(publicDir string) Option {
	return func(lp *litepage) error {
		lp.publicDir = publicDir
		return nil
	}
}

// By default a sitemap.xml is created mapping all pages of the static site. You can disable this
// behavior by specifying without a sitemap.
func WithoutSitemap() Option {
	return func(lp *litepage) error {
		lp.withSitemap = false
		return nil
	}
}

// Specify the base path of your site, if it is not the root of the domain. If set, all static assets and links should add the base as a prefix. The path should always start with a ”/” and not end with a trailing slash (otherwise an error will be returned).
func WithBasePath(basePath string) Option {
	return func(lp *litepage) error {
		if err := validate.IsValidBasePath(basePath); err != nil {
			return fmt.Errorf("base path is not valid '%s': %w", basePath, err)
		}
		lp.basePath = basePath
		return nil
	}
}

func (lp *litepage) Page(filePath string, handler func(w io.Writer)) error {
	err := validate.IsValidFilePath(filePath)
	if err != nil {
		return fmt.Errorf("error when validating file path '%s': %w", filePath, err)
	}
	exists := lp.pathMap[filePath]
	if exists {
		return fmt.Errorf("cannot add page '%s', it already exists", filePath)
	}

	// Because there is no native ordered map in Go, I opted to use a slice to preserve
	// page insertion order, and copy the file path into a map for o(1) lookup time.
	lp.pathMap[filePath] = true
	*lp.pages = append(*lp.pages, model.Page{Path: filePath, Handler: handler})

	return nil
}

func (lp *litepage) Serve(port string) error {
	sc := serve.Config{
		PublicDir:   lp.publicDir,
		Pages:       lp.pages,
		SiteDomain:  lp.siteDomain,
		BasePath:    lp.basePath,
		WithSitemap: lp.withSitemap,
	}
	server := serve.New(sc)
	return server.Serve(port)
}

func (lp *litepage) Build() error {
	bc := build.Config{
		DistDir:     lp.distDir,
		PublicDir:   lp.publicDir,
		Pages:       lp.pages,
		SiteDomain:  lp.siteDomain,
		BasePath:    lp.basePath,
		WithSitemap: lp.withSitemap,
	}
	builder := build.New(bc)
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
