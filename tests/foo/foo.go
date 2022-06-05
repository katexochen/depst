package foo

import (
	"github.com/katexochen/depst/tests/a"
	"github.com/katexochen/depst/tests/b"
	"github.com/katexochen/depst/tests/dependonaandb"
	"github.com/katexochen/depst/tests/dependonstd"
	"github.com/katexochen/depst/tests/nodependency"
)

var (
	Foo  = a.A + b.B + dependonaandb.Value + nodependency.N
	time = dependonstd.T
)
