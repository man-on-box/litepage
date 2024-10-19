package common

import "io"

type PageMap map[string]func(w io.Writer)

type Page struct {
	Path    string
	Handler func(w io.Writer)
}
