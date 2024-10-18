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
	PageMap     *common.PageMap
	SiteDomain  string
	WithSitemap bool
}

func New(publicDir string, pageMap *common.PageMap, siteDomain string, withSitemap bool) SiteServer {
	s := &siteServer{
		PublicDir:   publicDir,
		PageMap:     pageMap,
		SiteDomain:  siteDomain,
		WithSitemap: withSitemap,
	}
	return s
}

func (s *siteServer) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	var rootHandler http.HandlerFunc

	sortedPaths := common.SortPageMapByPath(s.PageMap)
	for _, path := range sortedPaths {
		fmt.Println("- serving page: ", path)

		pageHandler := (*s.PageMap)[path]
		handler := func(w http.ResponseWriter, r *http.Request) {
			pageHandler(w)
		}

		pathWithoutExt := strings.TrimSuffix(path, filepath.Ext(path))
		mux.HandleFunc(path, handler)
		mux.HandleFunc(pathWithoutExt, handler)

		if pathWithoutExt == "/index" {
			rootHandler = handler
		}
	}

	if s.WithSitemap {
		sitemap := common.BuildSitemap(s.SiteDomain, s.PageMap)
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
