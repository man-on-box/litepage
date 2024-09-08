package pages

import (
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/man-on-box/litepage/data"
	"github.com/man-on-box/litepage/util"
)

const distDir = "dist"

type Page struct {
	render   func(w io.Writer)
	filePath string
}

type Config struct {
	Domain string
}

type Pages struct {
	config Config
	data   *data.Data
	tmpl   *template.Template
	pages  *[]Page
}

func New(config Config) *Pages {
	p := Pages{
		config: config,
		data:   data.New(config.Domain),
	}
	return &p
}
func (p *Pages) init() {
	p.tmpl = parseTemplates()
	scaffoldDistDir(distDir)
	copyPublicDir(distDir)
}

func parseTemplates() *template.Template {
	patterns := []string{
		"./view/*.html",
		"./view/**/*.html",
	}

	tmpl := template.New("").Funcs(template.FuncMap{
		"version": func() string {
			return time.Now().Format("01021504")
		},
	})
	var err error

	for _, pattern := range patterns {
		tmpl, err = tmpl.ParseGlob(pattern)
		if err != nil {
			log.Fatalf("Error parsing templates: %v", err)
		}
	}

	return tmpl
}

func scaffoldDistDir(distDir string) {
	os.MkdirAll(distDir+"/articles", os.ModePerm)
}

func copyPublicDir(distDir string) {
	if err := util.CopyDir("public", distDir); err != nil {
		log.Fatalf("Could not copy public directory: %v", err)
	}
}
