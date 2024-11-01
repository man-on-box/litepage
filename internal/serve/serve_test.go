package serve_test

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/man-on-box/litepage/internal/common"
	"github.com/man-on-box/litepage/internal/serve"
	"github.com/stretchr/testify/assert"
)

func TestSiteServer(t *testing.T) {
	body := map[string]string{
		"index":               "<h1>Index Page</h1>",
		"foo":                 "<h1>Foo Page</h1>",
		"nested-index":        "<h1>Nested Index Page</h1>",
		"nested-foo":          "<h1>Nested Foo Page</h1>",
		"nested-nested-index": "<h1>Nested nested Index Page</h1>",
		"nested-nested-bar":   "<h1>Nested nested Bar Page</h1>",
		"text-file-body":      "example text response",
		"testfile":            "Hello from static text file",
		"sitemap":             `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://test.com/</loc></url><url><loc>https://test.com/foo</loc></url><url><loc>https://test.com/nested/</loc></url><url><loc>https://test.com/nested/foo</loc></url><url><loc>https://test.com/nested/nested/</loc></url><url><loc>https://test.com/nested/nested/bar</loc></url></urlset>`,
	}

	testPages := &[]common.Page{
		{
			Path: "/index.html",
			Handler: func(w io.Writer) {
				t := template.Must(template.New("").Parse(body["index"]))
				t.Execute(w, nil)
			},
		},
		{
			Path: "/foo.htm",
			Handler: func(w io.Writer) {
				t := template.Must(template.New("").Parse(body["foo"]))
				t.Execute(w, nil)
			},
		},
		{
			Path: "/nested/index.htm",
			Handler: func(w io.Writer) {
				t := template.Must(template.New("").Parse(body["nested-index"]))
				t.Execute(w, nil)
			},
		},
		{
			Path: "/nested/foo.htm",
			Handler: func(w io.Writer) {
				t := template.Must(template.New("").Parse(body["nested-foo"]))
				t.Execute(w, nil)
			},
		},
		{
			Path: "/nested/nested/index.html",
			Handler: func(w io.Writer) {
				t := template.Must(template.New("").Parse(body["nested-nested-index"]))
				t.Execute(w, nil)
			},
		},
		{
			Path: "/nested/nested/bar.html",
			Handler: func(w io.Writer) {
				t := template.Must(template.New("").Parse(body["nested-nested-bar"]))
				t.Execute(w, nil)
			},
		},
		{
			Path: "/textfile-endpoint.txt",
			Handler: func(w io.Writer) {
				w.Write([]byte(body["text-file-body"]))
			},
		},
	}

	tmpPublicDir, err := os.MkdirTemp("", "public")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpPublicDir)

	fmt.Printf("tmpPublicDir: %s\n", tmpPublicDir)

	testFilePath := tmpPublicDir + "/testfile.txt"
	err = os.WriteFile(testFilePath, []byte(body["testfile"]), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Can serve index page at '/'",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   body["index"],
		},
		{
			name:           "Can serve index page at '/index.html'",
			path:           "/index.html",
			expectedStatus: http.StatusOK,
			expectedBody:   body["index"],
		},
		{
			name:           "Can serve index page at '/index'",
			path:           "/index",
			expectedStatus: http.StatusOK,
			expectedBody:   body["index"],
		},
		{
			name:           "Can serve foo page at '/foo'",
			path:           "/foo",
			expectedStatus: http.StatusOK,
			expectedBody:   body["foo"],
		},
		{
			name:           "Can serve foo page at '/foo.htm'",
			path:           "/foo.htm",
			expectedStatus: http.StatusOK,
			expectedBody:   body["foo"],
		},
		{
			name:           "Can serve nested foo page at '/nested/foo.htm'",
			path:           "/nested/foo.htm",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-foo"],
		},
		{
			name:           "Can serve nested foo page at '/nested/foo'",
			path:           "/nested/foo",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-foo"],
		},
		{
			name:           "Can serve nested index page as htm file at '/nested/index.htm'",
			path:           "/nested/index.htm",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-index"],
		},
		{
			name:           "Can serve nested index page at '/nested/index'",
			path:           "/nested/index",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-index"],
		},
		{
			name:           "Can serve nested index page at '/nested', without the trailing slash",
			path:           "/nested",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-index"],
		},
		{
			name:           "Can serve nested index page at '/nested/', with the trailing slash",
			path:           "/nested/",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-index"],
		},
		{
			name:           "Can serve nested nested bar page with ext",
			path:           "/nested/nested/bar.html",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-nested-bar"],
		},
		{
			name:           "Can serve nested bar page without ext",
			path:           "/nested/nested/bar",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-nested-bar"],
		},
		{
			name:           "Can serve nested nested index page at '/nested/nested/index'",
			path:           "/nested/nested/index",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-nested-index"],
		},
		{
			name:           "Can serve nested nested index page at '/nested/nested', without the trailing slash",
			path:           "/nested/nested",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-nested-index"],
		},
		{
			name:           "Can serve nested nested index page at '/nested/nested/', with the trailing slash",
			path:           "/nested/nested/",
			expectedStatus: http.StatusOK,
			expectedBody:   body["nested-nested-index"],
		},
		{
			name:           "Returns 404 for non existent page",
			path:           "/nope",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
		{
			name:           "Returns 404 for non existent nested page",
			path:           "/nested/nope",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
		{
			name:           "Returns text body from .txt file",
			path:           "/textfile-endpoint.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   body["text-file-body"],
		},
		{
			name:           "Does not load text file without .txt extension",
			path:           "/textfile-endpoint",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
		{
			name:           "Returns test file from public dir",
			path:           "/testfile.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   body["testfile"],
		},
		{
			name:           "Returns expected sitemap",
			path:           "/sitemap.xml",
			expectedStatus: http.StatusOK,
			expectedBody:   body["sitemap"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serve.New(tmpPublicDir, testPages, "test.com", true)
			routes := s.SetupRoutes()
			server := httptest.NewServer(routes)
			defer server.Close()

			resp, err := http.Get(server.URL + tt.path)
			assert.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, string(body))
		})
	}
}
