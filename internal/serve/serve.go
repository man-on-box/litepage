package serve

import (
	"fmt"
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

type siteServer struct {
	PublicDir   string
	Pages       *[]common.Page
	SiteDomain  string
	WithSitemap bool
}

func New(publicDir string, pages *[]common.Page, siteDomain string, withSitemap bool) SiteServer {
	s := &siteServer{
		PublicDir:   publicDir,
		Pages:       pages,
		SiteDomain:  siteDomain,
		WithSitemap: withSitemap,
	}
	return s
}

func (s *siteServer) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	var rootHandler http.HandlerFunc

	for _, p := range *s.Pages {
		fmt.Println("- serving page: ", p.Path)

		handler := func(w http.ResponseWriter, r *http.Request) {
			p.Handler(w)
		}

		pathWithoutExt := strings.TrimSuffix(p.Path, filepath.Ext(p.Path))
		mux.HandleFunc(p.Path, handler)
		mux.HandleFunc(pathWithoutExt, handler)

		if pathWithoutExt == "/index" {
			rootHandler = handler
		}
	}

	if s.WithSitemap {
		sitemap := common.BuildSitemap(s.SiteDomain, s.Pages)
		mux.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(sitemap))
		})
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && rootHandler != nil {
			rootHandler(w, r)
		} else {
			http.ServeFile(w, r, s.PublicDir+r.URL.Path)
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
