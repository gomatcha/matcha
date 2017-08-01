package cmd

import (
	"go/build"
)

// func (ctxt *Context) ImportDir(dir string, mode ImportMode) (*Package, error)

func ImportAll(ctx *build.Context, paths []string, srcDir string, mode build.ImportMode) (map[string]*build.Package, error) {
	pkgs := map[string]*build.Package{}
	for _, i := range paths {
		if err := Import(ctx, i, srcDir, mode, pkgs); err != nil {
			return nil, err
		}
	}

	// fmt.Println(pkg)
	// fmt.Println("Dir", pkg.Dir)
	// fmt.Println("Name", pkg.Name)
	// fmt.Println("ImportComment", pkg.ImportComment)
	// fmt.Println("ImportPath", pkg.ImportPath)
	// fmt.Println("Root", pkg.Root)
	// fmt.Println("SrcRoot", pkg.SrcRoot)
	// fmt.Println("PkgRoot", pkg.PkgRoot)
	// fmt.Println("Goroot", pkg.Goroot)
	// fmt.Println("Imports", pkg.Imports)

	return pkgs, nil
}

func Import(ctx *build.Context, path, srcDir string, mode build.ImportMode, pkgs map[string]*build.Package) error {
	// fmt.Println(
	// Ignore C
	if path == "C" {
		return nil
	}

	pkg, err := ctx.Import(path, srcDir, mode)
	if err != nil {
		return err
	}

	// Bail if we have already added this package
	if _, ok := pkgs[pkg.Dir]; ok {
		return nil
	}
	pkgs[pkg.Dir] = pkg

	// Import dependencies
	for _, i := range pkg.Imports {
		if err := Import(ctx, i, pkg.Dir, mode, pkgs); err != nil {
			return err
		}
	}
	return nil
}

// Dir /Users/Overcyn/Workspace/go/src/gomatcha.io/matcha/examples
// Name examples
// ImportComment
// ImportPath gomatcha.io/matcha/examples
// Root /Users/Overcyn/Workspace/go/
// SrcRoot /Users/Overcyn/Workspace/go/src
// PkgRoot /Users/Overcyn/Workspace/go/pkg
// Goroot false
// Imports [gomatcha.io/matcha/examples/animate gomatcha.io/matcha/examples/complex gomatcha.io/matcha/examples/constraints gomatcha.io/matcha/examples/imageview gomatcha.io/matcha/examples/paint gomatcha.io/matcha/examples/screen gomatcha.io/matcha/examples/settings gomatcha.io/matcha/examples/stackscreen gomatcha.io/matcha/examples/table gomatcha.io/matcha/examples/tabscreen gomatcha.io/matcha/examples/textview gomatcha.io/matcha/examples/touch]
