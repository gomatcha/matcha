package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Init(flags *Flags) error {
	start := time.Now()

	// Get $GOPATH/pkg/gomobile
	gomobilepath, err := GoMobilePath()
	if err != nil {
		return err
	}
	if flags.ShouldPrint() {
		fmt.Fprintln(os.Stderr, "GOMOBILE="+gomobilepath)
	}

	// Delete $GOPATH/pkg/gomobile
	if err := RemoveAll(flags, gomobilepath); err != nil {
		return err
	}

	// Make $GOPATH/pkg/gomobile
	if err := Mkdir(flags, gomobilepath); err != nil {
		return err
	}

	// Make $GOPATH/pkg/gomobile/work...
	tmpdir, err := NewTmpDir(flags, gomobilepath)
	if err != nil {
		return err
	}
	defer RemoveAll(flags, tmpdir)

	// Install standard libraries for cross compilers.
	var env []string
	if env, err = DarwinArmEnv(flags); err != nil {
		return err
	}
	if err := InstallPkg(flags, tmpdir, "std", env); err != nil {
		return err
	}

	if env, err = DarwinArm64Env(flags); err != nil {
		return err
	}
	if err := InstallPkg(flags, tmpdir, "std", env); err != nil {
		return err
	}

	if env, err = Darwin386Env(flags); err != nil {
		return err
	}
	if err := InstallPkg(flags, tmpdir, "std", env, "-tags=ios"); err != nil {
		return err
	}

	if env, err = DarwinAmd64Env(flags); err != nil {
		return err
	}
	if err := InstallPkg(flags, tmpdir, "std", env, "-tags=ios"); err != nil {
		return err
	}

	// Write Go Version to $GOPATH/pkg/gomobile/version
	verpath := filepath.Join(gomobilepath, "version")
	if flags.ShouldPrint() {
		fmt.Fprintln(os.Stderr, "go version >", verpath)
	}
	if flags.ShouldRun() {
		goversion, err := GoVersion(flags)
		if err != nil {
			return nil
		}
		if err := ioutil.WriteFile(verpath, goversion, 0644); err != nil {
			return err
		}
	}

	// Timing
	if flags.BuildV {
		took := time.Since(start) / time.Second * time.Second
		fmt.Fprintf(os.Stderr, "Build took %s.\n", took)
	}
	fmt.Fprintf(os.Stderr, "Matcha initialized.\n")
	return nil
}

// Build package with properties.
func InstallPkg(f *Flags, temporarydir string, pkg string, env []string, args ...string) error {
	pd, err := PkgPath(env)
	if err != nil {
		return err
	}

	tOS, tArch := Getenv(env, "GOOS"), Getenv(env, "GOARCH")
	if tOS != "" && tArch != "" {
		if f.BuildV {
			fmt.Fprintf(os.Stderr, "\n# Installing %s for %s/%s.\n", pkg, tOS, tArch)
		}
		args = append(args, "-pkgdir="+pd)
	} else {
		if f.BuildV {
			fmt.Fprintf(os.Stderr, "\n# Installing %s.\n", pkg)
		}
	}

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
	return RunCmd(f, temporarydir, cmd)
}
