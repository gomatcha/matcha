// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func validateXcodeInstall(flags *Flags) error {
	err := _validateXcodeInstall(flags)
	if err != nil {
		fmt.Println(`Invalid or unsupported Xcode installation. See https://gomatcha.io/guide/installation/
for detailed instructions or set the --targets="android" flag to skip iOS builds.
`)
	}
	return err
}

func _validateXcodeInstall(flags *Flags) error {
	if _, err := LookPath(flags, "xcrun"); err != nil {
		return err
	}
	return nil
}

func ArchClang(goarch string) string {
	switch goarch {
	case "arm":
		return "armv7"
	case "arm64":
		return "arm64"
	case "386":
		return "i386"
	case "amd64":
		return "x86_64"
	default:
		panic(fmt.Sprintf("unknown GOARCH: %q", goarch))
	}
}

// Get clang path and clang flags (SDK Path).
func EnvClang(flags *Flags, sdkName string) (_clang, cflags string, err error) {
	// Get the clang path
	cmd := exec.Command("xcrun", "--sdk", sdkName, "--find", "clang")

	var clang string
	if flags.ShouldPrint() {
		PrintCmd(cmd)
	}
	if flags.ShouldRun() {
		out, err := cmd.Output()
		if err != nil {
			return "", "", fmt.Errorf("xcrun --find: %v\n%s", err, out)
		}
		clang = strings.TrimSpace(string(out))
	} else {
		clang = "clang-" + sdkName
	}

	// Get the SDK path
	cmd = exec.Command("xcrun", "--sdk", sdkName, "--show-sdk-path")
	var sdk string
	if flags.ShouldPrint() {
		PrintCmd(cmd)
	}
	if flags.ShouldRun() {
		out, err := cmd.Output()
		if err != nil {
			return "", "", fmt.Errorf("xcrun --show-sdk-path: %v\n%s", err, out)
		}
		sdk = strings.TrimSpace(string(out))
	} else {
		sdk = sdkName
	}

	return clang, "-isysroot " + sdk, nil
}

func DarwinArmEnv(f *Flags) ([]string, error) {
	clang, cflags, err := EnvClang(f, "iphoneos")
	if err != nil {
		return nil, err
	}
	return []string{
		"GOOS=darwin",
		"GOARCH=arm",
		"GOARM=7",
		"CC=" + clang,
		"CXX=" + clang,
		"CGO_CFLAGS=" + cflags + " -miphoneos-version-min=6.1 -arch " + ArchClang("arm"),
		"CGO_LDFLAGS=" + cflags + " -miphoneos-version-min=6.1 -arch " + ArchClang("arm"),
		"CGO_ENABLED=1",
	}, nil
}

func DarwinArm64Env(f *Flags) ([]string, error) {
	clang, cflags, err := EnvClang(f, "iphoneos")
	if err != nil {
		return nil, err
	}
	return []string{
		"GOOS=darwin",
		"GOARCH=arm64",
		"CC=" + clang,
		"CXX=" + clang,
		"CGO_CFLAGS=" + cflags + " -miphoneos-version-min=6.1 -arch " + ArchClang("arm64"),
		"CGO_LDFLAGS=" + cflags + " -miphoneos-version-min=6.1 -arch " + ArchClang("arm64"),
		"CGO_ENABLED=1",
	}, nil
}

func Darwin386Env(f *Flags) ([]string, error) {
	clang, cflags, err := EnvClang(f, "iphonesimulator")
	if err != nil {
		return nil, err
	}
	return []string{
		"GOOS=darwin",
		"GOARCH=386",
		"CC=" + clang,
		"CXX=" + clang,
		"CGO_CFLAGS=" + cflags + " -mios-simulator-version-min=6.1 -arch " + ArchClang("386"),
		"CGO_LDFLAGS=" + cflags + " -mios-simulator-version-min=6.1 -arch " + ArchClang("386"),
		"CGO_ENABLED=1",
	}, nil
}

func DarwinAmd64Env(f *Flags) ([]string, error) {
	clang, cflags, err := EnvClang(f, "iphonesimulator")
	if err != nil {
		return nil, err
	}
	return []string{
		"GOOS=darwin",
		"GOARCH=amd64",
		"CC=" + clang,
		"CXX=" + clang,
		"CGO_CFLAGS=" + cflags + " -mios-simulator-version-min=6.1 -arch x86_64",
		"CGO_LDFLAGS=" + cflags + " -mios-simulator-version-min=6.1 -arch x86_64",
		"CGO_ENABLED=1",
	}, nil
}
