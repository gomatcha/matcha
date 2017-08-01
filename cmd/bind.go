package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func Build(flags *Flags, args []string) error {
	iosDir, err := PackageDir(flags, "gomatcha.io/matcha/ios")
	if err != nil {
		return err
	}
	flags.BuildBinary = true
	flags.BuildO = iosDir
	return Bind(flags, args)
}

func Bind(flags *Flags, args []string) error {
	// Make $WORK.
	tempdir, err := NewTmpDir(flags, "")
	if err != nil {
		return err
	}
	if !flags.BuildWork {
		defer RemoveAll(flags, tempdir)
	}

	// Make $WORK/matcha-ios
	workOutputDir := filepath.Join(tempdir, "matcha-ios")
	if err := Mkdir(flags, workOutputDir); err != nil {
		return err
	}

	// Get $GOPATH/pkg/gomobile.
	gomobilepath, err := GoMobilePath()
	if err != nil {
		return err
	}

	// Get toolchain version.
	installedVersion, err := ReadFile(flags, filepath.Join(gomobilepath, "version"))
	if err != nil {
		return errors.New("toolchain partially installed, run `matcha init`")
	}

	// Get go version.
	goVersion, err := GoVersion(flags)
	if err != nil {
		return err
	}

	// Check toolchain matches go version.
	if !bytes.Equal(installedVersion, goVersion) && flags.ShouldRun() {
		return errors.New("toolchain out of date, run `matcha init`")
	}

	// Get current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create a build context.
	ctx := build.Default
	ctx.GOARCH = "arm"
	ctx.GOOS = "darwin"
	ctx.BuildTags = append(ctx.BuildTags, "ios")

	// Get import paths to be built.
	importPaths := []string{}
	srcDir := ""
	if len(args) == 0 {
		importPaths = append(importPaths, ".")
		srcDir = cwd
	} else {
		for _, i := range args {
			i = path.Clean(i)
			importPaths = append(importPaths, i)
		}
		srcDir = cwd
	}

	// Get packages to be built
	pkgs, err := ImportAll(&ctx, importPaths, srcDir, build.ImportComment)
	if err != nil {
		return err
	}

	// Check if any of the package is main.
	for _, pkg := range pkgs {
		if pkg.Name == "main" {
			return fmt.Errorf("binding 'main' package (%s) is not supported", pkg.ImportComment)
		}
	}

	genDir := filepath.Join(tempdir, "gen")
	binaryPath := filepath.Join(workOutputDir, "MatchaBridge", "MatchaBridge", "MatchaBridge.a")
	if err := Mkdir(flags, filepath.Dir(binaryPath)); err != nil {
		return err
	}

	// Build the "matcha/bridge" dir
	bridgeDir := filepath.Join(genDir, "src", "gomatcha.io", "bridge")
	if err := Mkdir(flags, bridgeDir); err != nil {
		return err
	}

	// Create the "main" go package, that references the other go packages
	mainPath := filepath.Join(tempdir, "src", "iosbin", "main.go")
	err = WriteFile(flags, mainPath, func(w io.Writer) error {
		format := fmt.Sprintf(BindFile, args[0]) // TODO(KD): Should this be args[0] or should it use the logic to generate pkgs
		_, err := w.Write([]byte(format))
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to create the binding package for iOS: %v", err)
	}

	// Get the supporting files
	cmdPath, err := PackageDir(flags, "gomatcha.io/matcha/cmd")
	if err != nil {
		return err
	}
	if err := CopyFile(flags, filepath.Join(bridgeDir, "matchaobjc.h"), filepath.Join(cmdPath, "matchaobjc.h.support")); err != nil {
		return err
	}
	if err := CopyFile(flags, filepath.Join(bridgeDir, "matchaobjc.m"), filepath.Join(cmdPath, "matchaobjc.m.support")); err != nil {
		return err
	}
	if err := CopyFile(flags, filepath.Join(bridgeDir, "matchaobjc.go"), filepath.Join(cmdPath, "matchaobjc.go.support")); err != nil {
		return err
	}
	if err := CopyFile(flags, filepath.Join(bridgeDir, "matchago.h"), filepath.Join(cmdPath, "matchago.h.support")); err != nil {
		return err
	}
	if err := CopyFile(flags, filepath.Join(bridgeDir, "matchago.m"), filepath.Join(cmdPath, "matchago.m.support")); err != nil {
		return err
	}
	if err := CopyFile(flags, filepath.Join(bridgeDir, "matchago.go"), filepath.Join(cmdPath, "matchago.go.support")); err != nil {
		return err
	}

	if !flags.BuildBinary {
		// Copy package's ios directory if it imports gomatcha.io/bridge.
		for _, pkg := range pkgs {
			importsBridge := false
			for _, i := range pkg.Imports {
				if i == "gomatcha.io/bridge" {
					importsBridge = true
					break
				}
			}

			if importsBridge {
				files, err := ioutil.ReadDir(pkg.Dir)
				if err != nil {
					continue
				}

				for _, i := range files {
					if i.IsDir() && i.Name() == "ios" {
						// Copy directory
						src := filepath.Join(pkg.Dir, "ios")
						dst := filepath.Join(workOutputDir)
						CopyDirContents(flags, dst, src)
					}
				}
			}
		}

		// Copy headers into Xcode project.
		if err = CopyFile(flags, filepath.Join(workOutputDir, "MatchaBridge", "MatchaBridge", "matchaobjc.h"), filepath.Join(cmdPath, "matchaobjc.h.support")); err != nil {
			return err
		}
		if err = CopyFile(flags, filepath.Join(workOutputDir, "MatchaBridge", "MatchaBridge", "matchago.h"), filepath.Join(cmdPath, "matchago.h.support")); err != nil {
			return err
		}
	}

	// Build platform binaries concurrently.
	matchaDarwinArmEnv, err := DarwinArmEnv(flags)
	if err != nil {
		return err
	}
	matchaDarwinArm64Env, err := DarwinArm64Env(flags)
	if err != nil {
		return err
	}
	matchaDarwin386Env, err := Darwin386Env(flags)
	if err != nil {
		return err
	}
	matchaDarwinAmd64Env, err := DarwinAmd64Env(flags)
	if err != nil {
		return err
	}

	type archPath struct {
		arch string
		path string
		err  error
	}
	archChan := make(chan archPath)
	for _, i := range [][]string{matchaDarwinArmEnv, matchaDarwinArm64Env, matchaDarwinAmd64Env, matchaDarwin386Env} {
		go func(env []string) {
			arch := Getenv(env, "GOARCH")
			env = append(env, "GOPATH="+genDir+string(filepath.ListSeparator)+os.Getenv("GOPATH"))
			path := filepath.Join(tempdir, "matcha-"+arch+".a")
			err := GoBuild(flags, mainPath, env, ctx, tempdir, "-buildmode=c-archive", "-o", path)
			archChan <- archPath{arch, path, err}
		}(i)
	}
	archs := []archPath{}
	for i := 0; i < 4; i++ {
		arch := <-archChan
		if arch.err != nil {
			return arch.err
		}
		archs = append(archs, arch)
	}

	// Lipo to build fat binary.
	cmd := exec.Command("xcrun", "lipo", "-create")
	for _, i := range archs {
		cmd.Args = append(cmd.Args, "-arch", ArchClang(i.arch), i.path)
	}
	cmd.Args = append(cmd.Args, "-o", binaryPath)
	if err := RunCmd(flags, tempdir, cmd); err != nil {
		return err
	}

	// Create output dir
	outputDir := flags.BuildO
	if outputDir == "" {
		outputDir = "Matcha-iOS"
	}

	if !flags.BuildBinary {
		if err := RemoveAll(flags, outputDir); err != nil {
			return err
		}

		// Copy output directory into place.
		if err := CopyDir(flags, outputDir, workOutputDir); err != nil {
			return err
		}
	} else {
		// Copy binary into place.
		if err := CopyFile(flags, filepath.Join(outputDir, "MatchaBridge", "MatchaBridge", "MatchaBridge.a"), binaryPath); err != nil {
			return err
		}
	}
	return nil
}

var BindFile = `
package main

import (
    _ "gomatcha.io/bridge"
    _ "%s"
)

import "C"

func main() {}
`
