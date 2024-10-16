package serve

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/man-on-box/litepage/pkg/types"
)

type SiteServer interface {
	Serve(port string) error
}

type siteServer struct {
	PublicDir string
	Pages     *[]types.Page
}

const defaultPort = "3000"

func New(publicDir string, pages *[]types.Page) SiteServer {
	s := &siteServer{
		PublicDir: publicDir,
		Pages:     pages,
	}
	return s
}

func (s *siteServer) Serve(port string) error {
	usePort := port
	if usePort == "" {
		usePort = defaultPort
	}
	fmt.Printf("LITEPAGE starting dev server on port %s...\n", usePort)

	var rootHandler http.HandlerFunc

	for _, page := range *s.Pages {
		fmt.Println("- serving page: ", page.FilePath)

		handler := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s: %s\n", r.Method, r.URL)
			page.Handler(w)
		}

		pathWithoutExt := strings.TrimSuffix(page.FilePath, filepath.Ext(page.FilePath))
		http.HandleFunc(page.FilePath, handler)
		http.HandleFunc(pathWithoutExt, handler)

		if pathWithoutExt == "/index" {
			rootHandler = handler
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && rootHandler != nil {
			rootHandler(w, r)
		} else {
			http.ServeFile(w, r, "public"+r.URL.Path)
		}
	})

	return http.ListenAndServe("localhost:"+usePort, nil)
}
