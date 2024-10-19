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
		"index":    "<h1>Index Page</h1>",
		"foo":      "<h1>Foo Page</h1>",
		"testfile": "Hello from text file",
		"sitemap":  `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://test.com/</loc></url><url><loc>https://test.com/nested/foo</loc></url></urlset>`,
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
			Path: "/nested/foo.htm",
			Handler: func(w io.Writer) {
				t := template.Must(template.New("").Parse(body["foo"]))
				t.Execute(w, nil)
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
			name:           "Can serve index page at '/index",
			path:           "/index",
			expectedStatus: http.StatusOK,
			expectedBody:   body["index"],
		},
		{
			name:           "Can serve nested foo page at '/nested/foo.htm",
			path:           "/nested/foo.htm",
			expectedStatus: http.StatusOK,
			expectedBody:   body["foo"],
		},
		{
			name:           "Can serve nested foo page at '/nested/foo",
			path:           "/nested/foo",
			expectedStatus: http.StatusOK,
			expectedBody:   body["foo"],
		},
		{
			name:           "Returns 404 for non existent page",
			path:           "/nope",
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
