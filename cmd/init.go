// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmd provides supporting functions for the matcha command line tool.
package cmd

import (
	"bytes"
	"path/filepath"
	"time"
)

func Init(f *Flags) error {
	start := time.Now()

	// Validate Go
	err := validateGoInstall(f)
	if err != nil {
		return err
	}

	// Parse targets
	targets := ParseTargets(f.BuildTargets)

	// Get $GOPATH/pkg/matcha directory
	matchaPkgPath, err := MatchaPkgPath(f)
	if err != nil {
		return err
	}

	// Delete $GOPATH/pkg/matcha
	if err := RemoveAll(f, matchaPkgPath); err != nil {
		return err
	}

	// Make $GOPATH/pkg/matcha
	if err := Mkdir(f, matchaPkgPath); err != nil {
		return err
	}

	// Make $WORK
	tmpdir, err := NewTmpDir(f, "")
	if err != nil {
		return err
	}
	defer RemoveAll(f, tmpdir)

	// Begin iOS
	if _, ok := targets["ios"]; ok {
		// Validate Xcode installation
		if err := validateXcodeInstall(f); err != nil {
			return err
		}

		// Install standard libraries for cross compilers.
		var env []string

		if _, ok := targets["ios/arm"]; ok {
			if env, err = DarwinArmEnv(f); err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env); err != nil {
				return err
			}
		}

		if _, ok := targets["ios/arm64"]; ok {
			if env, err = DarwinArm64Env(f); err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env); err != nil {
				return err
			}
		}

		if _, ok := targets["ios/386"]; ok {
			if env, err = Darwin386Env(f); err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env, "-tags=ios"); err != nil {
				return err
			}
		}

		if _, ok := targets["ios/amd64"]; ok {
			if env, err = DarwinAmd64Env(f); err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env, "-tags=ios"); err != nil {
				return err
			}
		}
	}

	// Begin android
	if _, ok := targets["android"]; ok {
		// Validate Android installation
		if err := ValidateAndroidInstall(f); err != nil {
			return err
		}

		// Install standard libraries for cross compilers.
		if _, ok := targets["android/arm"]; ok {
			env, err := AndroidEnv(f, "arm")
			if err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env); err != nil {
				return err
			}
		}

		if _, ok := targets["android/arm64"]; ok {
			env, err := AndroidEnv(f, "arm64")
			if err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env); err != nil {
				return err
			}
		}

		if _, ok := targets["android/386"]; ok {
			env, err := AndroidEnv(f, "386")
			if err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env); err != nil {
				return err
			}
		}

		if _, ok := targets["android/amd64"]; ok {
			env, err := AndroidEnv(f, "amd64")
			if err != nil {
				return err
			}
			if err := InstallPkg(f, matchaPkgPath, tmpdir, "std", env); err != nil {
				return err
			}
		}
	}

	// Write Go Version to $GOPATH/pkg/matcha/version
	goversion, err := GoVersion(f)
	if err != nil {
		return nil
	}
	verpath := filepath.Join(matchaPkgPath, "version")
	if err := WriteFile(f, verpath, bytes.NewReader(goversion)); err != nil {
		return err
	}

	// Timing
	if f.BuildV {
		took := time.Since(start) / time.Second * time.Second
		f.Logger.Printf("Build took %s.\n", took)
	}
	f.Logger.Printf("Matcha initialized.\n")
	return nil
}
