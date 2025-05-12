package serve

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/man-on-box/litepage/internal/model"
	"github.com/man-on-box/litepage/internal/sitemap"
)

const defaultPort = "3000"

//go:embed 404.html
var notFoundHTML string

var notFoundTmpl = template.Must(template.New("404").Parse(notFoundHTML))

type SiteServer interface {
	Serve(port string) error
	SetupRoutes() http.Handler
}

type Config struct {
	PublicDir   string
	Pages       *[]model.Page
	SiteDomain  string
	BasePath    string
	WithSitemap bool
}

type siteServer struct {
	Config Config
}

func New(config Config) SiteServer {
	s := &siteServer{
		Config: config,
	}
	return s
}

func (s *siteServer) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	var rootHandler func(w io.Writer)
	registeredPaths := map[string]bool{}

	registerHandler := func(path string, handler func(w io.Writer)) {
		registeredPaths[path] = true

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if registeredPaths[r.URL.Path] {
				handler(w)
			} else {
				s.customNotFound(w, r)
			}
		})
	}

	for _, p := range *s.Config.Pages {
		fullPath := s.Config.BasePath + p.Path
		pageHandler := func(w io.Writer) {
			log.Printf("[%d]: %s", http.StatusOK, fullPath)
			p.Handler(w)

		}
		registerHandler(fullPath, pageHandler)

		fileExt := filepath.Ext(p.Path)
		isNotHTML := fileExt != ".html" && fileExt != ".htm"
		if isNotHTML {
			// do not have to register any other handlers
			continue
		}

		pathWithoutExt := strings.TrimSuffix(fullPath, fileExt)
		registerHandler(pathWithoutExt, pageHandler)

		if strings.HasSuffix(pathWithoutExt, "/index") {
			if pathWithoutExt == s.Config.BasePath+"/index" {
				rootHandler = pageHandler
			} else {
				trimmedPath := strings.TrimSuffix(pathWithoutExt, "/index")
				registerHandler(trimmedPath, pageHandler)
				registerHandler(trimmedPath+"/", pageHandler)
			}
		}
	}

	if s.Config.WithSitemap {
		smap := sitemap.BuildSitemap(s.Config.SiteDomain, s.Config.Pages)
		mux.HandleFunc(s.Config.BasePath+"/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(smap))
		})
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == s.Config.BasePath+"/" && rootHandler != nil {
			rootHandler(w)
			return
		}
		if len(s.Config.BasePath) > 0 && !strings.HasPrefix(r.URL.Path, s.Config.BasePath) {
			s.customNotFound(w, r)
			return
		}
		s.serveFile(w, r)
	})
	return mux
}

func (s *siteServer) Serve(port string) error {
	usePort := port
	if usePort == "" {
		usePort = defaultPort
	}
	fmt.Printf("LITEPAGE starting dev server at http://localhost:%s...\n", usePort)

	mux := s.SetupRoutes()
	return http.ListenAndServe("localhost:"+usePort, mux)
}

func (s *siteServer) serveFile(w http.ResponseWriter, r *http.Request) {
	staticPath := strings.TrimPrefix(r.URL.Path, s.Config.BasePath)
	requestedFile := filepath.Join(s.Config.PublicDir, filepath.Clean(staticPath))
	if _, err := os.Stat(requestedFile); os.IsNotExist(err) {
		s.customNotFound(w, r)
	} else {
		http.ServeFile(w, r, s.Config.PublicDir+staticPath)
	}
}

func (s *siteServer) customNotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%d]: %s", http.StatusNotFound, r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	showBasePathError := len(s.Config.BasePath) > 0 && !strings.HasPrefix(r.URL.Path, s.Config.BasePath)
	notFoundTmpl.Execute(w, struct {
		ShowBasePathError bool
		BasePath          string
		UrlPath           string
	}{ShowBasePathError: showBasePathError, BasePath: s.Config.BasePath, UrlPath: r.URL.Path})
}
