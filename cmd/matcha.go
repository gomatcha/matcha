// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Flags struct {
	BuildN       bool   // print commands but don't run
	BuildX       bool   // print commands
	BuildV       bool   // print package names
	BuildWork    bool   // use working directory
	BuildGcflags string // -gcflags
	BuildLdflags string // -ldflags
	BuildO       string // output path
	BuildBinary  bool
	BuildTargets string
}

func (f *Flags) ShouldPrint() bool {
	return f.BuildN || f.BuildX
}

func (f *Flags) ShouldRun() bool {
	return !f.BuildN
}

func FindEnv(env []string, key string) string {
	prefix := key + "="
	for _, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			return kv[len(prefix):]
		}
	}
	return ""
}

// $GOPATH/pkg/matcha
func MatchaPkgPath(f *Flags) (string, error) {
	gopaths := filepath.SplitList(GoEnv(f, "GOPATH"))
	gomobilepath := ""
	for _, p := range gopaths {
		gomobilepath = filepath.Join(p, "pkg", "matcha")
		if _, err := os.Stat(gomobilepath); err == nil {
			break
		}
	}
	if gomobilepath == "" {
		if len(gopaths) == 0 {
			return "", fmt.Errorf("$GOPATH does not exist")
		} else {
			return filepath.Join(gopaths[0], "pkg", "matcha"), nil
		}
	}
	return gomobilepath, nil
}

// $GOPATH/pkg/matcha/pkg_darwin_arm64
func PkgPath(f *Flags, env []string) (string, error) {
	gomobilepath, err := MatchaPkgPath(f)
	if err != nil {
		return "", err
	}
	return gomobilepath + "/pkg_" + FindEnv(env, "GOOS") + "_" + FindEnv(env, "GOARCH"), nil
}

// Returns the go enviromental variable for name.
func GoEnv(f *Flags, name string) string {
	if val := GetEnv(f, name); val != "" {
		return val
	}

	cmd := exec.Command("go", "env", name)
	if f.ShouldPrint() {
		PrintCmd(cmd)
	}
	if f.ShouldRun() {
		val, err := cmd.Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(val))
	}
	return "$" + name
}

func GoVersion(f *Flags) ([]byte, error) {
	cmd := exec.Command("go", "version")
	if f.ShouldPrint() {
		PrintCmd(cmd)
	}
	if f.ShouldRun() {
		goVer, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("'go version' failed: %v, %s", err, goVer)
		}
		if bytes.HasPrefix(goVer, []byte("go version go1.4")) || bytes.HasPrefix(goVer, []byte("go version go1.5")) || bytes.HasPrefix(goVer, []byte("go version go1.6")) {
			return nil, errors.New("Go 1.7 or newer is required")
		}
		return goVer, nil
	}
	return []byte("go version goX.X.X"), nil
}

func GoBuild(f *Flags, src string, env []string, ctx build.Context, tmpdir string, args ...string) error {
	return GoCmd(f, "build", []string{src}, env, ctx, tmpdir, args...)
}

func GoInstall(f *Flags, srcs []string, env []string, ctx build.Context, tmpdir string, args ...string) error {
	return GoCmd(f, "install", srcs, env, ctx, tmpdir, args...)
}

func GoCmd(f *Flags, subcmd string, srcs []string, env []string, ctx build.Context, tmpdir string, args ...string) error {
	pkgPath, err := PkgPath(f, env)
	if err != nil {
		return err
	}

	if !IsDir(f, pkgPath) {
		return fmt.Errorf("Matcha not initialized for this target. Missing directory at %v.", pkgPath)
	}

	cmd := exec.Command("go", subcmd, "-pkgdir="+pkgPath)
	if len(ctx.BuildTags) > 0 {
		cmd.Args = append(cmd.Args, "-tags", strings.Join(ctx.BuildTags, " "))
	}
	if f.BuildV {
		cmd.Args = append(cmd.Args, "-v")
	}
	// if subcmd != "install" && f.BuildI {
	// 	cmd.Args = append(cmd.Args, "-i")
	// }
	if f.BuildX {
		cmd.Args = append(cmd.Args, "-x")
	}
	if f.BuildGcflags != "" {
		cmd.Args = append(cmd.Args, "-gcflags", f.BuildGcflags)
	}
	if f.BuildLdflags != "" {
		cmd.Args = append(cmd.Args, "-ldflags", f.BuildLdflags)
	}
	if f.BuildWork {
		cmd.Args = append(cmd.Args, "-work")
	}
	cmd.Args = append(cmd.Args, args...)
	cmd.Args = append(cmd.Args, srcs...)
	cmd.Env = append([]string{}, env...)
	return RunCmd(f, tmpdir, cmd)
}
