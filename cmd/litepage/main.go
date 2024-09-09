package main

import (
	"github.com/man-on-box/litepage/pages"
)

func main() {
	pages := pages.New(pages.Config{
		Domain: "yourdomain.com",
	})
	pages.Create()
}
