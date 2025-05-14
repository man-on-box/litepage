package sitemap_test

import (
	"testing"

	"github.com/man-on-box/litepage/internal/model"
	"github.com/man-on-box/litepage/internal/sitemap"
	"github.com/stretchr/testify/assert"
)

func TestSiteMap(t *testing.T) {
	testPages := &[]model.Page{
		{
			Path: "/index.html",
		},
		{
			Path: "/foo.htm",
		},
		{
			Path: "/nested/index.htm",
		},
		{
			Path: "/nested/foo.htm",
		},
		{
			Path: "/nested/nested/index.html",
		},
		{
			Path: "/nested/nested/bar.html",
		},
		{
			Path: "/a-text-file.txt",
		},
	}

	t.Run("creates expected sitemap from pages", func(t *testing.T) {
		smap := sitemap.Build("test.com", "", testPages)

		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://test.com/</loc></url><url><loc>https://test.com/foo</loc></url><url><loc>https://test.com/nested/</loc></url><url><loc>https://test.com/nested/foo</loc></url><url><loc>https://test.com/nested/nested/</loc></url><url><loc>https://test.com/nested/nested/bar</loc></url></urlset>`

		assert.Equal(t, expected, smap)
	})

	t.Run("creates expected sitemap from pages including when base path is specified", func(t *testing.T) {
		smap := sitemap.Build("test.com", "/test", testPages)

		expected := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://test.com/test/</loc></url><url><loc>https://test.com/test/foo</loc></url><url><loc>https://test.com/test/nested/</loc></url><url><loc>https://test.com/test/nested/foo</loc></url><url><loc>https://test.com/test/nested/nested/</loc></url><url><loc>https://test.com/test/nested/nested/bar</loc></url></urlset>`

		assert.Equal(t, expected, smap)
	})
}
