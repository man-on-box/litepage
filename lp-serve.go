package litepage

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func (lp *LitePage) Serve(port string) error {
	fmt.Printf("LITEPAGE starting dev server on port %s...\n", port)
	var rootHandler http.HandlerFunc

	for _, page := range *lp.pages {
		fmt.Println("- serving page: ", page.filePath)

		handler := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s: %s\n", r.Method, r.URL)
			page.handler(w)
		}

		pathWithoutExt := strings.TrimSuffix(page.filePath, filepath.Ext(page.filePath))
		http.HandleFunc(page.filePath, handler)
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

	return http.ListenAndServe("localhost:"+port, nil)
}
