package depst

import (
	"fmt"
	"testing"
)

func BenchmarkLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		im := make(importMap)
		im.load("github.com/katexochen/depst/tests/foo")
	}
}

func TestLoad(t *testing.T) {
	im := make(importMap)
	err := im.load("github.com/katexochen/depst/tests/...")
	if err != nil {
		t.Fatal(err)
	}
	printImportMap(im)
}

func printImportMap(im importMap) {
	for key, imps := range im {
		fmt.Println(key)
		for imp := range imps {
			fmt.Printf("\t%s\n", imp)
		}
	}
}
