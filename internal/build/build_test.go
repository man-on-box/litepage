package build_test

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"testing"

	"github.com/man-on-box/litepage/internal/build"
	"github.com/man-on-box/litepage/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestSiteBuilder(t *testing.T) {
	body := map[string]string{
		"index":          "<h1>Index Page</h1>",
		"foo":            "<h1>Foo Page</h1>",
		"text-file-body": "example text response",
		"testfile":       "Hello from text file",
		"sitemap":        `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://test.com/</loc></url><url><loc>https://test.com/nested/foo</loc></url><url><loc>https://test.com/zzz</loc></url><url><loc>https://test.com/aaa</loc></url></urlset>`,
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
		{
			Path: "/text-file-path.txt",
			Handler: func(w io.Writer) {
				w.Write([]byte(body["text-file-body"]))
			},
		},
		{
			Path:    "/zzz.html",
			Handler: func(w io.Writer) {},
		},
		{
			Path:    "/aaa.html",
			Handler: func(w io.Writer) {},
		},
	}

	tmpDistDir, err := os.MkdirTemp("", "dist")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDistDir)

	tmpPublicDir, err := os.MkdirTemp("", "public")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpPublicDir)

	testFilePath := tmpPublicDir + "/testfile.txt"
	err = os.WriteFile(testFilePath, []byte(body["testfile"]), 0644)
	assert.NoError(t, err)

	b := build.New(tmpDistDir, tmpPublicDir, testPages, "test.com", true)
	err = b.Build()
	assert.NoError(t, err)

	tests := []struct {
		path            string
		expectedContent string
	}{
		{
			path:            "/index.html",
			expectedContent: body["index"],
		},
		{
			path:            "/nested/foo.htm",
			expectedContent: body["foo"],
		},
		{
			path:            "/text-file-path.txt",
			expectedContent: body["text-file-body"],
		},
		{
			path:            "/testfile.txt",
			expectedContent: body["testfile"],
		},
		{
			path:            "/sitemap.xml",
			expectedContent: body["sitemap"],
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("File '%s' exists and contains expected content", tt.path), func(t *testing.T) {
			content, err := os.ReadFile(tmpDistDir + tt.path)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedContent, string(content))
		})
	}
}
