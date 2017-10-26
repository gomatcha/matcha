// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	return pkgs, nil
}

func Import(ctx *build.Context, path, srcDir string, mode build.ImportMode, pkgs map[string]*build.Package) error {
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
