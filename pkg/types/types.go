package types

import "io"

type Page struct {
	FilePath string
	Handler  PageHandler
}

type PageHandler func(w io.Writer)
