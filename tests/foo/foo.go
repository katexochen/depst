package foo

import (
	"github.com/katexochen/depst/tests/a"
	"github.com/katexochen/depst/tests/b"
	"github.com/katexochen/depst/tests/dependonaandb"
	"github.com/katexochen/depst/tests/dependonstd"
)

var (
	Foo  = a.A + b.B + dependonaandb.Value
	time = dependonstd.T
)
