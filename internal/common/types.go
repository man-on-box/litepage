package common

import "io"

type PageMap map[string]func(w io.Writer)
