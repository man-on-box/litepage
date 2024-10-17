package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/man-on-box/litepage"
)

func main() {
	lp, err := litepage.New("example.dev")
	if err != nil {
		log.Fatalf("Could not create app: %v", err)
	}

	lp.Page("/index.html", handleHomepage())

	err = lp.BuildOrServe()
	if err != nil {
		log.Fatalf("Could not run app: %v", err)
	}
}

func handleHomepage() func(w io.Writer) {
	data := struct {
		Title     string `json:"title"`
		Header    string `json:"header"`
		Subheader string `json:"subheader"`
		DocsUrl   string `json:"docsUrl"`
	}{}
	parseJSONFile("./content/homepage.json", &data)

	t := template.Must(tmpl.ParseFiles("./view/base.html", "./view/home.html"))

	return func(w io.Writer) {
		err := t.ExecuteTemplate(w, "base", data)
		if err != nil {
			log.Fatalf("Could not execute template: %v", err)
		}
	}

}

var tmpl = template.New("").Funcs(template.FuncMap{
	"version": func() string {
		return time.Now().Format("01021504")
	},
})

func parseJSONFile(file string, data any) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}

	err = json.Unmarshal(contents, &data)
	if err != nil {
		log.Fatalf("Could not parse JSON: %v", err)
	}
}
