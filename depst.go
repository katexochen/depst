package depst

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"
)

/*
buildflags
external deps break
std deps break
*/

type Depst struct {
	importMap importMap
}

func New() *Depst {
	return &Depst{make(importMap)}
}

func (d Depst) Eval(r rule) error {
	var err error
	r.a, err = expand(r.a)
	if err != nil {
		return err
	}
	if err := d.importMap.load(r.a...); err != nil {
		return err
	}
	r.b, err = expand(r.b)
	if err != nil {
		return err
	}
	if err := d.importMap.load(r.b...); err != nil {
		return err
	}
	return eval(d.importMap, r)
}

func Packages(pkgs ...string) pkgPaths {
	return pkgs
}

type pkgPaths []string

func (p pkgPaths) ShouldNotDependOn(pkgs ...string) rule {
	return rule{
		a: p,
		b: pkgs,
	}
}

func (p pkgPaths) ShouldNotDependDirectlyOn(pkgs ...string) rule {
	return rule{
		a:        p,
		b:        pkgs,
		directly: true,
	}
}

func (p pkgPaths) ShouldDependOn(pkgs ...string) rule {
	return rule{
		a:      p,
		b:      pkgs,
		should: true,
	}
}

func (p pkgPaths) ShouldDependDirectlyOn(pkgs ...string) rule {
	return rule{
		a:        p,
		b:        pkgs,
		directly: true,
		should:   true,
	}
}

type rule struct {
	a        pkgPaths
	b        pkgPaths
	should   bool
	directly bool
}

func Eval(r rule) error {
	im := make(importMap)
	var err error
	r.a, err = expand(r.a)
	if err != nil {
		return err
	}
	if err := im.load(r.a...); err != nil {
		return err
	}
	r.b, err = expand(r.b)
	if err != nil {
		return err
	}
	if err := im.load(r.b...); err != nil {
		return err
	}

	return eval(im, r)
}

func eval(im importMap, r rule) error {
	for _, a := range r.a {
		for _, b := range r.b {
			if r.directly {
				if im.dependsDirectlyOn(a, b) != r.should {
					return fmt.Errorf("%s depends directly on %s", a, b)
				}
				return nil
			}

			if im.dependsOn(a, b) != r.should {
				return fmt.Errorf("%s depends on %s", a, b)
			}
		}
	}
	return nil
}

type importMap map[string]map[string]struct{} // from package => to package

func (i importMap) load(pkgs ...string) error {
	config := &packages.Config{Mode: packages.NeedName | packages.NeedImports}

	loadedPkgs, err := packages.Load(config, pkgs...)
	if err != nil {
		return err
	}
	// fmt.Println(loadedPkgs)

	var toLoad []string
	for _, pkg := range loadedPkgs {
		if _, ok := i[pkg.PkgPath]; ok {
			continue
		}

		i[pkg.PkgPath] = make(map[string]struct{})
		for _, imp := range pkg.Imports {
			i[pkg.PkgPath][imp.PkgPath] = struct{}{}

			if _, ok := i[imp.PkgPath]; !ok {
				toLoad = append(toLoad, imp.PkgPath)
			}
		}
	}

	if len(toLoad) > 0 {
		return i.load(toLoad...)
	}
	return nil
}

func (i importMap) dependsDirectlyOn(a, b string) bool {
	_, ok := i[a][b]
	return ok
}

func (i importMap) dependsOn(a, b string) bool {
	if i.dependsDirectlyOn(a, b) {
		return true
	}
	aDeps := i[a]
	for dep := range aDeps {
		if dep == b {
			return true
		}
		if i.dependsOn(dep, b) {
			return true
		}
	}
	return false
}

func expand(pkgs pkgPaths) (pkgPaths, error) {
	cfg := packages.Config{Mode: packages.NeedName}

	for i, pkg := range pkgs {
		if !strings.Contains(pkg, "...") {
			continue
		}

		expanded, err := packages.Load(&cfg, pkg)
		if err != nil {
			return nil, err
		}

		pkgs[i] = expanded[0].PkgPath

		for _, pkg := range expanded[1:] {
			pkgs = append(pkgs, pkg.PkgPath)
		}
	}

	return pkgs, nil
}
