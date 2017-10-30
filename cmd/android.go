// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"archive/zip"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	javacTargetVer = "1.7"
	minAndroidAPI  = 15
	manifestHeader = `Manifest-Version: 1.0
Created-By: 1.0 (Go)

`
)

const (
	missingAndroidHomeEnvVar  = "$ANDROID_HOME enviromental variable is unset and does not point to an Android SDK. The SDK is often at ~/Library/Android/sdk on macOS and ~/Android/Sdk on Linux."
	missingAndroidHome        = "$ANDROID_HOME enviromental variable does not point to an Android SDK. The SDK is often at ~/Library/Android/sdk on macOS and ~/Android/Sdk on Linux."
	missingNDK                = "NDK was not found at $ANDROID_HOME/ndk-bundle. NDK can be installed in Android Studio > SDK Manager."
	missingAndroidPlatformDir = "Android SDK platform directory was not found at $ANDROID_HOME/ndk-bundle/platforms."
	missingAndroidPlatform    = "Android SDK platform with minimum API level of 15 was not found in $ANDROID_HOME/ndk-bundle/platforms. SDK platforms can be installed in Android Studio > SDK Manager."
	missingJavac              = "javac was not found in $PATH."
)

func ValidateAndroidInstall(f *Flags) error {
	err := validateAndroidInstall(f)
	if err != nil {
		fmt.Println(`Invalid or unsupported Android installation. See https://gomatcha.io/guide/installation/
for detailed instructions or set the --targets="ios" flag to skip Android builds.
`)
	}
	return err
}

func validateAndroidInstall(f *Flags) error {
	if _, err := AndroidPlatformPath(f); err != nil {
		return err
	}
	if _, err := NDKPath(f); err != nil {
		return err
	}
	if _, err := LookPath(f, "javac"); err != nil {
		return fmt.Errorf(missingJavac)
	}
	return nil
}

func AndroidSDKPath(f *Flags) (string, error) {
	path := GetEnv(f, "ANDROID_HOME")
	if path == "" {
		return "", fmt.Errorf(missingAndroidHomeEnvVar)
	}

	if !IsDir(f, path) {
		return "", fmt.Errorf(missingAndroidHome)
	}
	return path, nil
}

// AndroidPlatformPath returns an android SDK platform directory under ANDROID_HOME.
// If there are multiple platforms that satisfy the minimum version requirement
// AndroidPlatformPath returns the latest one among them.
func AndroidPlatformPath(f *Flags) (string, error) {
	androidHome, err := AndroidSDKPath(f)
	if err != nil {
		return "", err
	}

	platformsDir := filepath.Join(androidHome, "platforms")
	if !IsDir(f, platformsDir) {
		return "", fmt.Errorf(missingAndroidPlatformDir)
	}

	platformsDirNames, err := ReadDirNames(f, platformsDir)
	if err != nil {
		return "", err
	}
	if !f.ShouldRun() {
		platformsDirNames = []string{"android-21"}
	}

	var apiPath string
	var apiVer int
	for _, i := range platformsDirNames {
		verStr := strings.TrimPrefix(i, "android-")
		if i == verStr {
			continue
		}

		ver, err := strconv.Atoi(verStr)
		if err != nil || ver < minAndroidAPI || ver < apiVer {
			continue
		}

		p := filepath.Join(platformsDir, i)
		if !IsFile(f, filepath.Join(p, "android.jar")) {
			continue
		}

		apiPath = p
		apiVer = ver
	}

	if apiVer == 0 {
		return "", fmt.Errorf(missingAndroidPlatform)
	}
	return apiPath, nil
}

func NDKPath(f *Flags) (string, error) {
	path, err := AndroidSDKPath(f)
	if err != nil {
		return "", err
	}

	path = filepath.Join(path, "ndk-bundle")
	if !IsDir(f, path) {
		return "", fmt.Errorf(missingNDK)
	}
	return path, nil
}

func AndroidEnv(f *Flags, goarch string) ([]string, error) {
	tc, err := toolchainForArch(f, goarch)
	if err != nil {
		return nil, err
	}
	flags := fmt.Sprintf("-target %s -gcc-toolchain %s", tc.clangTriple, tc.gccToolchain())
	cflags := fmt.Sprintf("%s --sysroot %s -isystem %s -D__ANDROID_API__=%s", flags, tc.csysroot(), tc.isystem(), tc.api)
	ldflags := fmt.Sprintf("%s --sysroot %s", flags, tc.ldsysroot())
	env := []string{
		"GOOS=android",
		"GOARCH=" + goarch,
		"CC=" + tc.clangPath(),
		"CXX=" + tc.clangppPath(),
		"CGO_CFLAGS=" + cflags,
		"CGO_CPPFLAGS=" + cflags,
		"CGO_LDFLAGS=" + ldflags,
		"CGO_ENABLED=1",
	}
	if goarch == "arm" {
		env = append(env, "GOARM=7")
	}
	return env, nil
}

// Emulate the flags in the clang wrapper scripts generated
// by make_standalone_toolchain.py
// https://android.googlesource.com/platform/ndk/+/ndk-release-r16/docs/UnifiedHeaders.md
// https://developer.android.com/ndk/guides/standalone_toolchain.html#c_stl_support
// http://zwyuan.github.io/2015/12/22/three-ways-to-use-android-ndk-cross-compiler/
type ndkToolchain struct {
	arch        string
	api         string
	gcc         string
	triple      string
	clangTriple string

	ndkRoot string
	hostTag string
}

func toolchainForArch(f *Flags, goarch string) (*ndkToolchain, error) {
	m := map[string]*ndkToolchain{
		"arm": &ndkToolchain{
			arch:        "arm",
			api:         "15",
			gcc:         "arm-linux-androideabi-4.9",
			triple:      "arm-linux-androideabi",
			clangTriple: "armv7a-none-linux-androideabi",
		},
		"arm64": &ndkToolchain{
			arch:        "arm64",
			api:         "21",
			gcc:         "aarch64-linux-android-4.9",
			triple:      "aarch64-linux-android",
			clangTriple: "aarch64-none-linux-android",
		},
		"386": &ndkToolchain{
			arch:        "x86",
			api:         "15",
			gcc:         "x86-4.9",
			triple:      "i686-linux-android",
			clangTriple: "i686-none-linux-android",
		},
		"amd64": &ndkToolchain{
			arch:        "x86_64",
			api:         "21",
			gcc:         "x86_64-4.9",
			triple:      "x86_64-linux-android",
			clangTriple: "x86_64-none-linux-android",
		},
	}
	toolchain, ok := m[goarch]
	if !ok {
		return nil, fmt.Errorf("toolchainForArch(): Unknown arch %v", goarch)
	}

	ndkRoot, err := NDKPath(f)
	if err != nil {
		return nil, err
	}
	toolchain.ndkRoot = ndkRoot

	hostTag, err := ndkHostTag()
	if err != nil {
		return nil, err
	}
	toolchain.hostTag = hostTag
	return toolchain, nil
}

func (tc *ndkToolchain) gccToolchain() string {
	return filepath.Join(tc.ndkRoot, "toolchains", tc.gcc, "prebuilt", tc.hostTag)
}

func (tc *ndkToolchain) clangPath() string {
	return filepath.Join(tc.ndkRoot, "toolchains", "llvm", "prebuilt", tc.hostTag, "bin", "clang")
}

func (tc *ndkToolchain) clangppPath() string {
	return filepath.Join(tc.ndkRoot, "toolchains", "llvm", "prebuilt", tc.hostTag, "bin", "clang++")
}

func (tc *ndkToolchain) isystem() string {
	return filepath.Join(tc.ndkRoot, "sysroot", "usr", "include", tc.triple)
}

func (tc *ndkToolchain) csysroot() string {
	return filepath.Join(tc.ndkRoot, "sysroot")
}

func (tc *ndkToolchain) ldsysroot() string {
	return filepath.Join(tc.ndkRoot, "platforms", "android-"+tc.api, "arch-"+tc.arch)
}

func GetAndroidABI(arch string) string {
	switch arch {
	case "arm":
		return "armeabi-v7a"
	case "arm64":
		return "arm64-v8a"
	case "386":
		return "x86"
	case "amd64":
		return "x86_64"
	}
	return ""
}

func ndkHostTag() (string, error) {
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		return "windows", nil
	} else {
		var arch string
		switch runtime.GOARCH {
		case "386":
			arch = "x86"
		case "amd64":
			arch = "x86_64"
		default:
			return "", fmt.Errorf("ndkHostTag(): Unsupported GOARCH %v", runtime.GOARCH)
		}
		return runtime.GOOS + "-" + arch, nil
	}
}

// AAR is the format for the binary distribution of an Android Library Project
// and it is a ZIP archive with extension .aar.
// http://tools.android.com/tech-docs/new-build-system/aar-format
//
// These entries are directly at the root of the archive.
//
//  AndroidManifest.xml (mandatory)
//  classes.jar (mandatory)
//  assets/ (optional)
//  jni/<abi>/libgojni.so
//  R.txt (mandatory)
//  res/ (mandatory)
//  libs/*.jar (optional, not relevant)
//  proguard.txt (optional)
//  lint.jar (optional, not relevant)
//  aidl (optional, not relevant)
//
// javac and jar commands are needed to build classes.jar.
func BuildAAR(f *Flags, androidDir string, pkgs []*build.Package, androidArchs []string, tmpdir string, aarPath string) (err error) {
	if !f.ShouldRun() { // TODO(KD):
		return nil
	}

	var out io.Writer = ioutil.Discard
	if !f.BuildN {
		f, err := os.Create(aarPath)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}()
		out = f
	}

	aarw := zip.NewWriter(out)
	aarwcreate := func(name string) (io.Writer, error) {
		if f.BuildV {
			f.Logger.Printf("aar: %s\n", name)
		}
		return aarw.Create(name)
	}
	w, err := aarwcreate("AndroidManifest.xml")
	if err != nil {
		return err
	}
	const manifestFmt = `<manifest xmlns:android="http://schemas.android.com/apk/res/android" package=%q>
<uses-sdk android:minSdkVersion="%d"/></manifest>`
	fmt.Fprintf(w, manifestFmt, "go."+pkgs[0].Name+".gojni", minAndroidAPI)

	w, err = aarwcreate("proguard.txt")
	if err != nil {
		return err
	}
	fmt.Fprintln(w, `-keep class go.** { *; }`)

	w, err = aarwcreate("classes.jar")
	if err != nil {
		return err
	}
	src := filepath.Join(androidDir, "src/main/java")
	if err := BuildJar(f, w, src, tmpdir); err != nil {
		return err
	}

	// Assets....
	// files := map[string]string{}
	// for _, pkg := range pkgs {
	// 	assetsDir := filepath.Join(pkg.Dir, "assets")
	// 	assetsDirExists := false
	// 	if fi, err := os.Stat(assetsDir); err == nil {
	// 		assetsDirExists = fi.IsDir()
	// 	} else if !os.IsNotExist(err) {
	// 		return err
	// 	}

	// 	if assetsDirExists {
	// 		err := filepath.Walk(
	// 			assetsDir, func(path string, info os.FileInfo, err error) error {
	// 				if err != nil {
	// 					return err
	// 				}
	// 				if info.IsDir() {
	// 					return nil
	// 				}
	// 				f, err := os.Open(path)
	// 				if err != nil {
	// 					return err
	// 				}
	// 				defer f.Close()
	// 				name := "assets/" + path[len(assetsDir)+1:]
	// 				if orig, exists := files[name]; exists {
	// 					return fmt.Errorf("package %s asset name conflict: %s already added from package %s",
	// 						pkg.ImportPath, name, orig)
	// 				}
	// 				files[name] = pkg.ImportPath
	// 				w, err := aarwcreate(name)
	// 				if err != nil {
	// 					return nil
	// 				}
	// 				_, err = io.Copy(w, f)
	// 				return err
	// 			})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	for _, arch := range androidArchs {
		lib := GetAndroidABI(arch) + "/libgojni.so"
		w, err = aarwcreate("jni/" + lib)
		if err != nil {
			return err
		}
		if !f.BuildN {
			r, err := os.Open(filepath.Join(androidDir, "src/main/jniLibs/"+lib))
			if err != nil {
				return err
			}
			defer r.Close()
			if _, err := io.Copy(w, r); err != nil {
				return err
			}
		}
	}

	// TODO(hyangah): do we need to use aapt to create R.txt?
	w, err = aarwcreate("R.txt")
	if err != nil {
		return err
	}

	w, err = aarwcreate("res/")
	if err != nil {
		return err
	}

	return aarw.Close()
}

func BuildJar(f *Flags, w io.Writer, srcDir string, tmpdir string) error {
	var srcFiles []string
	if !f.ShouldRun() {
		srcFiles = []string{"*.java"}
	} else {
		err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".java" {
				srcFiles = append(srcFiles, filepath.Join(".", path[len(srcDir):]))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	dst := filepath.Join(tmpdir, "javac-output")
	if err := Mkdir(f, dst); err != nil {
		return err
	}

	bClspath, err := bootClasspath(f)
	if err != nil {
		return err
	}

	args := []string{
		"-d", dst,
		"-source", javacTargetVer,
		"-target", javacTargetVer,
		"-bootclasspath", bClspath,
		// "-classpath", bindClasspath
	}
	args = append(args, srcFiles...)

	javac := exec.Command("javac", args...)
	javac.Dir = srcDir
	if err := RunCmd(f, tmpdir, javac); err != nil {
		return err
	}

	// fmt.Println("javac", args)
	// if buildX {
	// KD: printcmd("jar c -C %s .", dst)
	// }
	if !f.ShouldRun() {
		return nil
	}
	jarw := zip.NewWriter(w)
	jarwcreate := func(name string) (io.Writer, error) {
		if f.BuildV {
			f.Logger.Printf("jar: %s\n", name)
		}
		return jarw.Create(name)
	}
	manifestFile, err := jarwcreate("META-INF/MANIFEST.MF")
	if err != nil {
		return err
	}
	fmt.Fprintf(manifestFile, manifestHeader)

	err = filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		out, err := jarwcreate(filepath.ToSlash(path[len(dst)+1:]))
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		_, err = io.Copy(out, in)
		return err
	})
	if err != nil {
		return err
	}
	return jarw.Close()
}

func bootClasspath(f *Flags) (string, error) {
	// bindBootClasspath := "" // KD: command parameter
	// if bindBootClasspath != "" {
	// 	return bindBootClasspath, nil
	// }
	apiPath, err := AndroidPlatformPath(f)
	if err != nil {
		return "", err
	}
	return filepath.Join(apiPath, "android.jar"), nil
}
