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

func BenchmarkEval_ShouldNotDependOn_Nodependency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Eval(Packages("github.com/katexochen/depst/tests/...").ShouldNotDependOn("github.com/katexochen/depst/tests/nodependency"))
	}
}

func TestEval(t *testing.T) {
	testCases := map[string]struct {
		rule    rule
		testNeg bool
		wantErr bool
	}{
		"fail: dependona should not directly depend on a": {
			rule: Packages("github.com/katexochen/depst/tests/dependona").
				ShouldNotDependDirectlyOn("github.com/katexochen/depst/tests/a"),
			wantErr: true,
			testNeg: true,
		},
		"fail: indirectona should not depend on a": {
			rule: Packages("github.com/katexochen/depst/tests/indirectona").
				ShouldNotDependOn("github.com/katexochen/depst/tests/a"),
			wantErr: true,
			testNeg: true,
		},
		"fail: indirectona should directly depend on a": {
			rule: Packages("github.com/katexochen/depst/tests/indirectona").
				ShouldDependDirectlyOn("github.com/katexochen/depst/tests/a"),
			wantErr: true,
			testNeg: true,
		},
		"success: no one should directly depend on nodependency": {
			rule: Packages("github.com/katexochen/depst/tests/...").
				ShouldNotDependOn("github.com/katexochen/depst/tests/nodependency"),
		},
		// "success: no one except indirectona should directly depend on dependona": {
		// 	rule: Packages("github.com/katexochen/depst/tests/...").
		// 		Except("github.com/katexochen/depst/tests/indirectona").
		// 		ShouldNotDirectlyDependOn("github.com/katexochen/depst/tests/dependona"),
		// },
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := Eval(tc.rule); (err != nil) != tc.wantErr {
				t.Errorf("Eval(%v) = %v, want %v", tc.rule, err, tc.wantErr)
			}
			if !tc.testNeg {
				return
			}
			tc.rule.should = !tc.rule.should
			if err := Eval(tc.rule); (err != nil) == tc.wantErr {
				t.Errorf("Negation: Eval(%v) = %v, want %v", tc.rule, err, tc.wantErr)
			}
		})
	}
}

func TestDepsEval(t *testing.T) {
	testCases := map[string]struct {
		rule    rule
		testNeg bool
		wantErr bool
	}{
		"fail: dependona should not directly depend on a": {
			rule: Packages("github.com/katexochen/depst/tests/dependona").
				ShouldNotDependDirectlyOn("github.com/katexochen/depst/tests/a"),
			wantErr: true,
			testNeg: true,
		},
		"fail: indirectona should not depend on a": {
			rule: Packages("github.com/katexochen/depst/tests/indirectona").
				ShouldNotDependOn("github.com/katexochen/depst/tests/a"),
			wantErr: true,
			testNeg: true,
		},
		"fail: indirectona should directly depend on a": {
			rule: Packages("github.com/katexochen/depst/tests/indirectona").
				ShouldDependDirectlyOn("github.com/katexochen/depst/tests/a"),
			wantErr: true,
			testNeg: true,
		},
		"success: no one should directly depend on nodependency": {
			rule: Packages("github.com/katexochen/depst/tests/...").
				ShouldNotDependOn("github.com/katexochen/depst/tests/nodependency"),
		},
		// "success: no one except indirectona should directly depend on dependona": {
		// 	rule: Packages("github.com/katexochen/depst/tests/...").
		// 		Except("github.com/katexochen/depst/tests/indirectona").
		// 		ShouldNotDirectlyDependOn("github.com/katexochen/depst/tests/dependona"),
		// },
	}

	depst := New()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := depst.Eval(tc.rule); (err != nil) != tc.wantErr {
				t.Errorf("Eval(%v) = %v, want %v", tc.rule, err, tc.wantErr)
			}
			if !tc.testNeg {
				return
			}
			tc.rule.should = !tc.rule.should
			if err := depst.Eval(tc.rule); (err != nil) == tc.wantErr {
				t.Errorf("Negation: Eval(%v) = %v, want %v", tc.rule, err, tc.wantErr)
			}
		})
	}
}

func TestExpand(t *testing.T) {
	expanded, err := expand(Packages("github.com/katexochen/depst/tests/..."))
	if err != nil {
		t.Errorf("expand returned error: %v", err)
	}
	if len(expanded) != 9 {
		t.Errorf("expand = %v has length %d, want 9", expanded, len(expanded))
	}
}
