// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func ParseTargets(a string) map[string]struct{} {
	targetsSlice := strings.Fields(a)
	if len(targetsSlice) == 0 {
		if runtime.GOOS == "darwin" {
			targetsSlice = []string{"android", "ios"}
		} else {
			targetsSlice = []string{"android"}
		}
	}
	targets := map[string]struct{}{}
	for _, i := range targetsSlice {
		switch i {
		case "android":
			targets["android"] = struct{}{}
			targets["android/arm"] = struct{}{}
			targets["android/arm64"] = struct{}{}
			// targets["android/386"] = struct{}{}
			// targets["android/amd64"] = struct{}{}
		case "android/arm", "android/arm64", "android/386", "android/amd64":
			targets["android"] = struct{}{}
			targets[i] = struct{}{}
		case "ios":
			targets["ios"] = struct{}{}
			targets["ios/arm"] = struct{}{}
			targets["ios/arm64"] = struct{}{}
			// targets["ios/386"] = struct{}{}
			targets["ios/amd64"] = struct{}{}
		case "ios/arm", "ios/arm64", "ios/386", "ios/amd64":
			targets["ios"] = struct{}{}
			targets[i] = struct{}{}
		}
	}
	return targets
}

func Build(flags *Flags, args []string) error {
	iosDir, err := PackageDir(flags, "gomatcha.io/matcha")
	if err != nil {
		return err
	}
	// flags.BuildBinary = true
	flags.BuildO = iosDir
	return Bind(flags, args)
}

func Bind(flags *Flags, args []string) error {
	targets := ParseTargets(flags.BuildTargets)

	// Make $WORK.
	tempdir, err := NewTmpDir(flags, "")
	if err != nil {
		return err
	}
	if !flags.BuildWork {
		defer RemoveAll(flags, tempdir)
	}

	// Get $GOPATH/pkg/matcha.
	matchaPkgPath, err := MatchaPkgPath(flags)
	if err != nil {
		return err
	}

	// Get toolchain version.
	installedVersion, err := ReadFile(flags, filepath.Join(matchaPkgPath, "version"))
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
	cwd, err := Getwd(flags)
	if err != nil {
		return err
	}

	// Create a build context.
	ctx := build.Default
	ctx.GOARCH = "arm"
	ctx.GOOS = "darwin"
	ctx.BuildTags = append(ctx.BuildTags, "matcha")

	// Get import paths to be built.
	importPaths := []string{}
	if len(args) == 0 {
		importPaths = append(importPaths, ".")
	} else {
		for _, i := range args {
			i = path.Clean(i)
			importPaths = append(importPaths, i)
		}
	}

	// Get packages to be built
	pkgs, err := ImportAll(flags, &ctx, importPaths, cwd, build.ImportComment)
	if err != nil {
		return err
	}

	// Check if any of the package is main.
	for _, pkg := range pkgs {
		if pkg.Name == "main" {
			return fmt.Errorf("binding 'main' package (%s) is not supported", pkg.ImportComment)
		}
	}

	// Get the supporting files
	bridgePath, err := PackageDir(flags, "gomatcha.io/matcha/bridge")
	if err != nil {
		return err
	}

	// Begin iOS
	if _, ok := targets["ios"]; ok {
		// Validate Xcode installation
		if err := validateXcodeInstall(flags); err != nil {
			return err
		}

		// Build the "matcha/bridge" dir
		gopathDir := filepath.Join(tempdir, "IOS-GOPATH")

		// Make $WORK/matcha-ios
		workOutputDir := filepath.Join(tempdir, "matcha-ios")
		if err := Mkdir(flags, workOutputDir); err != nil {
			return err
		}

		// Make binary output dir
		binaryPath := filepath.Join(workOutputDir, "MatchaBridge", "MatchaBridge", "MatchaBridge.a")
		if err := Mkdir(flags, filepath.Dir(binaryPath)); err != nil {
			return err
		}

		// Create the "main" go package, that references the other go packages
		mainPath := filepath.Join(tempdir, "src", "iosbin", "main.go")
		err = WriteFile(flags, mainPath, strings.NewReader(fmt.Sprintf(BindFile, args[0]))) // TODO(KD): Should this be args[0] or should it use the logic to generate pkgs
		if err != nil {
			return fmt.Errorf("failed to create the binding package for iOS: %v", err)
		}

		// if !flags.BuildBinary {
		// 	// Copy package's ios directory if it imports gomatcha.io/bridge.
		// 	for _, pkg := range pkgs {
		// 		importsBridge := false
		// 		for _, i := range pkg.Imports {
		// 			if i == "gomatcha.io/bridge" {
		// 				importsBridge = true
		// 				break
		// 			}
		// 		}

		// 		if importsBridge {
		// 			files, err := ioutil.ReadDir(pkg.Dir)
		// 			if err != nil {
		// 				continue
		// 			}

		// 			for _, i := range files {
		// 				if i.IsDir() && i.Name() == "ios" {
		// 					// Copy directory
		// 					src := filepath.Join(pkg.Dir, "ios")
		// 					dst := filepath.Join(workOutputDir)
		// 					CopyDirContents(flags, dst, src)
		// 				}
		// 			}
		// 		}
		// 	}
		// }

		// Build platform binaries concurrently.
		envs := [][]string{}
		if _, ok := targets["ios/arm"]; ok {
			env, err := DarwinArmEnv(flags)
			if err != nil {
				return err
			}
			envs = append(envs, env)
		}
		if _, ok := targets["ios/arm64"]; ok {
			env, err := DarwinArm64Env(flags)
			if err != nil {
				return err
			}
			envs = append(envs, env)
		}
		if _, ok := targets["ios/386"]; ok {
			env, err := Darwin386Env(flags)
			if err != nil {
				return err
			}
			envs = append(envs, env)
		}
		if _, ok := targets["ios/amd64"]; ok {
			env, err := DarwinAmd64Env(flags)
			if err != nil {
				return err
			}
			envs = append(envs, env)
		}

		type archPath struct {
			arch string
			path string
			err  error
		}
		archs := []archPath{}
		archChan := make(chan archPath)
		for _, i := range envs {
			go func(env []string) {
				arch := FindEnv(env, "GOARCH")
				env = append(env, "GOPATH="+gopathDir+string(filepath.ListSeparator)+GoEnv(flags, "GOPATH"))
				path := filepath.Join(tempdir, "matcha-"+arch+".a")
				err := GoBuild(flags, []string{mainPath}, env, []string{"matcha"}, matchaPkgPath, tempdir, "-buildmode=c-archive", "-o", path)
				archChan <- archPath{arch, path, err}
			}(i)

			if !flags.Threaded {
				arch := <-archChan
				if arch.err != nil {
					return arch.err
				}
				archs = append(archs, arch)
			}
		}
		if flags.Threaded {
			for i := 0; i < len(envs); i++ {
				arch := <-archChan
				if arch.err != nil {
					return arch.err
				}
				archs = append(archs, arch)
			}
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

		// if !flags.BuildBinary {
		// 	if err := RemoveAll(flags, outputDir); err != nil {
		// 		return err
		// 	}

		// 	// Copy output directory into place.
		// 	if err := CopyDir(flags, outputDir, workOutputDir); err != nil {
		// 		return err
		// 	}
		// } else {

		// Copy binary into place.
		if err := CopyFile(flags, filepath.Join(outputDir, "ios", "MatchaBridge", "MatchaBridge", "MatchaBridge.a"), binaryPath); err != nil {
			return err
		}

		// }
	}
	if _, ok := targets["android"]; ok {
		// Validate Android installation
		if err := ValidateAndroidInstall(flags); err != nil {
			return err
		}

		// Build the "matcha/bridge" dir
		gopathDir := filepath.Join(tempdir, "ANDROID-GOPATH")

		androidArchs := []string{}
		if _, ok := targets["android/arm"]; ok {
			androidArchs = append(androidArchs, "arm")
		}
		if _, ok := targets["android/arm64"]; ok {
			androidArchs = append(androidArchs, "arm64")
		}
		if _, ok := targets["android/386"]; ok {
			androidArchs = append(androidArchs, "386")
		}
		if _, ok := targets["android/amd64"]; ok {
			androidArchs = append(androidArchs, "amd64")
		}

		androidDir := filepath.Join(tempdir, "android")
		mainPath := filepath.Join(tempdir, "androidlib/main.go")

		err = WriteFile(flags, mainPath, strings.NewReader(fmt.Sprintf(BindFile, args[0]))) // TODO(KD): Should this be args[0] or should it use the logic to generate pkgs
		if err != nil {
			return fmt.Errorf("failed to create the main package for android: %v", err)
		}

		javaDir2 := filepath.Join(androidDir, "src", "main", "java", "io", "gomatcha", "bridge")
		if err := Mkdir(flags, javaDir2); err != nil {
			return err
		}
		if err := CopyFile(flags, filepath.Join(javaDir2, "GoValue.java"), filepath.Join(bridgePath, "java-GoValue.java")); err != nil {
			return err
		}
		if err := CopyFile(flags, filepath.Join(javaDir2, "Bridge.java"), filepath.Join(bridgePath, "java-Bridge.java")); err != nil {
			return err
		}
		if err := CopyFile(flags, filepath.Join(javaDir2, "Tracker.java"), filepath.Join(bridgePath, "java-Tracker.java")); err != nil {
			return err
		}

		// Make $WORK/matcha-android
		workOutputDir := filepath.Join(tempdir, "matcha-android")
		if err := Mkdir(flags, workOutputDir); err != nil {
			return err
		}

		// Make aar output file.
		aarDirPath := filepath.Join(workOutputDir, "MatchaBridge")
		aarPath := filepath.Join(workOutputDir, "MatchaBridge", "matchabridge.aar")
		if err := Mkdir(flags, aarDirPath); err != nil {
			return err
		}

		// Generate binding code and java source code only when processing the first package.
		for _, arch := range androidArchs {
			env, err := AndroidEnv(flags, arch)
			if err != nil {
				return err
			}
			env = append(env, "GOPATH="+gopathDir+string(filepath.ListSeparator)+GoEnv(flags, "GOPATH"))

			err = GoBuild(flags,
				[]string{mainPath},
				env,
				[]string{"matcha"},
				matchaPkgPath,
				tempdir,
				"-buildmode=c-shared",
				"-o="+filepath.Join(androidDir, "src/main/jniLibs/"+GetAndroidABI(arch)+"/libgojni.so"),
			)
			if err != nil {
				return err
			}
		}

		if err := BuildAAR(flags, androidDir, pkgs, androidArchs, tempdir, aarPath); err != nil {
			return err
		}

		// Create output dir
		outputDir := flags.BuildO
		if outputDir == "" {
			outputDir = "Matcha-iOS"
		}

		// Copy binary into place.
		if err := CopyFile(flags, filepath.Join(outputDir, "android", "matchabridge.aar"), aarPath); err != nil {
			return err
		}
	}
	return nil
}

var BindFile = `
package main

import (
    _ "gomatcha.io/matcha/bridge"
    _ "%s"
)

import "C"

func main() {}
`
