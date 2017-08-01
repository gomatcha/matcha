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
}

func (f *Flags) ShouldPrint() bool {
	return f.BuildN || f.BuildX
}

func (f *Flags) ShouldRun() bool {
	return !f.BuildN
}

func XcodeAvailable() bool {
	_, err := exec.LookPath("xcrun")
	return err == nil
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
	if !XcodeAvailable() {
		return "", "", errors.New("Xcode not available")
	}

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

func Getenv(env []string, key string) string {
	prefix := key + "="
	for _, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			return kv[len(prefix):]
		}
	}
	return ""
}

// $GOPATH/pkg/gomobile
func GoMobilePath() (string, error) {
	gopaths := filepath.SplitList(GoEnv("GOPATH"))
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

// $GOPATH/pkg/gomobile/pkg_darwin_arm64
func PkgPath(env []string) (string, error) {
	gomobilepath, err := GoMobilePath()
	if err != nil {
		return "", err
	}
	return gomobilepath + "/pkg_" + Getenv(env, "GOOS") + "_" + Getenv(env, "GOARCH"), nil
}

// Returns the go enviromental variable for name.
func GoEnv(name string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	val, err := exec.Command("go", "env", name).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(val))
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
	pd, err := PkgPath(env)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", subcmd, "-pkgdir="+pd)
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
