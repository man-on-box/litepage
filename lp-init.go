package litepage

import (
	"errors"
)

type Option func(*LitePage)

func New(domain string, options ...Option) (*LitePage, error) {
	config := &config{
		siteDomain:  domain,
		distDir:     "dist",
		publicDir:   "public",
		withSitemap: true,
	}
	flags := &flags{
		serve: "lp-serve",
		port:  "lp-port",
	}
	lp := &LitePage{
		config: config,
		flags:  flags,
		pages:  &[]Page{},
	}

	for _, opt := range options {
		opt(lp)
	}

	if lp.config.siteDomain == "" {
		return nil, errors.New("site domain is required, please provide a domain like 'catpics.com'")
	}

	return lp, nil
}

func WithDistDir(distDis string) Option {
	return func(lp *LitePage) {
		lp.config.distDir = distDis
	}
}

func WithPublicDir(publicDir string) Option {
	return func(lp *LitePage) {
		lp.config.publicDir = publicDir
	}
}

func WithoutSitemap() Option {
	return func(lp *LitePage) {
		lp.config.withSitemap = false
	}
}

func WithCustomFlags(serveFlag string, portFlag string) Option {
	return func(lp *LitePage) {
		lp.flags = &flags{
			serve: serveFlag,
			port:  portFlag,
		}
	}
}
