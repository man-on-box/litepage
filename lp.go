package litepage

import (
	"flag"
	"io"
)

type Litepage struct {
	config *config
	flags  *flags
	pages  *[]Page
}

type config struct {
	siteDomain  string
	distDir     string
	publicDir   string
	withSitemap bool
}

type flags struct {
	serve string
	port  string
}

type Page struct {
	filePath string
	handler  func(w io.Writer)
}

func (lp *Litepage) Page(filePath string, handler func(w io.Writer)) {
	*lp.pages = append(*lp.pages, Page{filePath: filePath, handler: handler})
}

// Run the litepage app. By default, build static site to the dist directory.
// If -lp-serve flag is passed, it will instead serve the static site on port
// 3000, or on port specified if -lp-port param is passed.
func (lp *Litepage) Run() error {
	serve := flag.Bool(lp.flags.serve, false, "LITEPAGE will serve pages via dev server instead of creating the static site")
	port := flag.String(lp.flags.port, "3000", "Port the LITEPAGE dev server will start on (default: 3000)")
	flag.Parse()

	if *serve {
		return lp.Serve(*port)
	} else {
		return lp.Build()
	}
}
