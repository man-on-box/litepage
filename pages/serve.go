package pages

import (
	"fmt"
	"log"
	"net/http"
)

const devPort = ":3001"

func (p *Pages) serve() {
	fmt.Println("Serving pages...")

	for _, page := range *p.pages {
		path := fileToUrlPath(page.filePath)
		fmt.Println("Serving page: ", path)
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s: %s\n", r.Method, path)
			if path == "/" {
				if r.URL.Path == "/" {
					page.render(w)
				} else {
					http.ServeFile(w, r, "public"+r.URL.Path)
				}
			} else {
				page.render(w)
			}

		})
	}

	err := http.ListenAndServe("localhost"+devPort, nil)
	if err != nil {
		log.Fatalf("Could not start dev server: %v", err)
	}
}
