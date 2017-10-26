// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type Flags struct {
	Logger       *log.Logger
	Threaded     bool
	disablePrint bool
	BuildN       bool   // print commands but don't run
	BuildX       bool   // print commands
	BuildV       bool   // print package names. Verbose.
	BuildWork    bool   // use working directory
	BuildGcflags string // -gcflags
	BuildLdflags string // -ldflags
	BuildO       string // output path
	BuildBinary  bool
	BuildTargets string // targets
}

func (f *Flags) ShouldPrint() bool {
	return (f.BuildN || f.BuildX) && !f.disablePrint
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
	p := ""
	for _, i := range gopaths {
		p = filepath.Join(i, "pkg", "matcha")
		if IsDir(f, p) {
			break
		}
	}
	if p == "" {
		if len(gopaths) == 0 {
			return "", fmt.Errorf("$GOPATH does not exist")
		} else {
			return filepath.Join(gopaths[0], "pkg", "matcha"), nil
		}
	}
	return p, nil
}

// $GOPATH/pkg/matcha/pkg_darwin_arm64
func PkgPath(f *Flags, matchaPkgPath string, env []string) (string, error) {
	tOS, tArch := FindEnv(env, "GOOS"), FindEnv(env, "GOARCH")
	if tOS == "" || tArch == "" {
		return "", fmt.Errorf("PkgPath(): Missing GOOS or GOARCH", tOS, tArch)
	}

	return matchaPkgPath + "/pkg_" + tOS + "_" + tArch, nil
}

// Returns the go enviromental variable for name.
func GoEnv(f *Flags, name string) string {
	if val := GetEnv(f, name); val != "" {
		return val
	}

	cmd := exec.Command("go", "env", name)
	out, err := OutputCmd(f, []byte("$"+name), "", cmd)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func GoVersion(f *Flags) ([]byte, error) {
	cmd := exec.Command("go", "version")
	goVer, err := OutputCmd(f, []byte("go version goX.X.X x/x"), "", cmd)
	if err != nil {
		return nil, err
	}
	if f.ShouldRun() {
		if bytes.HasPrefix(goVer, []byte("go version go1.4")) || bytes.HasPrefix(goVer, []byte("go version go1.5")) || bytes.HasPrefix(goVer, []byte("go version go1.6")) {
			return nil, errors.New("Go 1.7 or newer is required")
		}
	}
	return goVer, nil
}

func GoBuild(f *Flags, srcs []string, env []string, buildTags []string, matchaPkgPath, tmpdir string, args ...string) error {
	pkgPath, err := PkgPath(f, matchaPkgPath, env)
	if err != nil {
		return err
	}

	if !IsDir(f, pkgPath) {
		return fmt.Errorf("Matcha not initialized for this target. Missing directory at %v.", pkgPath)
	}

	cmd := exec.Command("go", "build", "-pkgdir="+pkgPath)
	if len(buildTags) > 0 {
		cmd.Args = append(cmd.Args, "-tags", strings.Join(buildTags, " "))
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

// Build package with properties.
func InstallPkg(f *Flags, matchaPkgPath, temp string, pkg string, env []string, args ...string) error {
	pkgPath, err := PkgPath(f, matchaPkgPath, env)
	if err != nil {
		return err
	}
	args = append(args, "-pkgdir="+pkgPath)

	cmd := exec.Command("go", "install")
	cmd.Args = append(cmd.Args, args...)
	if f.BuildV {
		cmd.Args = append(cmd.Args, "-v")
	}
	if f.BuildX {
		cmd.Args = append(cmd.Args, "-x")
	}
	if f.BuildWork {
		cmd.Args = append(cmd.Args, "-work")
	}
	cmd.Args = append(cmd.Args, pkg)
	cmd.Env = append([]string{}, env...)
	return RunCmd(f, temp, cmd)
}
