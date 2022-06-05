package depst

import (
	"fmt"

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

type importMap map[string]map[string]struct{} // from package => to package

func (i importMap) load(pkgs ...string) error {
	config := &packages.Config{Mode: packages.NeedName | packages.NeedImports}

	loadedPkgs, err := packages.Load(config, pkgs...)
	if err != nil {
		return err
	}
	fmt.Println(loadedPkgs)

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

func isExternal(pkg string) bool {
	return pkg[0] != '.'
}
