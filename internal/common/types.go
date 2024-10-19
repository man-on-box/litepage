package common

import "io"

type Page struct {
	Path    string
	Handler func(w io.Writer)
}
