package serve

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/man-on-box/litepage/internal/common"
)

const defaultPort = "3000"

type SiteServer interface {
	Serve(port string) error
	SetupRoutes() http.Handler
}

type Config struct {
	PublicDir   string
	Pages       *[]common.Page
	SiteDomain  string
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
	registeredPaths := map[string]any{}

	registerHandler := func(path string, handler func(w io.Writer)) {
		registeredPaths[path] = struct{}{}

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			_, exists := registeredPaths[r.URL.Path]
			if exists {
				handler(w)
			} else {
				http.NotFound(w, r)
			}
		})
	}

	for _, p := range *s.Config.Pages {
		fmt.Println("- serving page: ", p.Path)
		registerHandler(p.Path, p.Handler)

		fileExt := filepath.Ext(p.Path)
		isNotHTML := fileExt != ".html" && fileExt != ".htm"
		if isNotHTML {
			// do not have to register any other handlers
			continue
		}

		pathWithoutExt := strings.TrimSuffix(p.Path, fileExt)
		registerHandler(pathWithoutExt, p.Handler)

		if strings.HasSuffix(pathWithoutExt, "/index") {
			if pathWithoutExt == "/index" {
				rootHandler = p.Handler
			} else {
				path := strings.TrimSuffix(pathWithoutExt, "/index")
				registerHandler(path, p.Handler)
				registerHandler(path+"/", p.Handler)
			}
		}
	}

	if s.Config.WithSitemap {
		sitemap := common.BuildSitemap(s.Config.SiteDomain, s.Config.Pages)
		mux.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(sitemap))
		})
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && rootHandler != nil {
			rootHandler(w)
		} else {
			http.ServeFile(w, r, s.Config.PublicDir+r.URL.Path)
		}
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
