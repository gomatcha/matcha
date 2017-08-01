package cmd

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func PrintCmd(cmd *exec.Cmd) {
	fmt.Fprintln(os.Stderr, strings.Join(cmd.Args, " "))
}

func RunCmd(f *Flags, tmpdir string, cmd *exec.Cmd) error {
	if f.ShouldPrint() {
		dir := ""
		if cmd.Dir != "" {
			dir = "PWD=" + cmd.Dir + " "
		}
		env := strings.Join(cmd.Env, " ")
		if env != "" {
			env += " "
		}
		fmt.Fprintln(os.Stderr, dir, env, strings.Join(cmd.Args, " "))
	}

	buf := new(bytes.Buffer)
	buf.WriteByte('\n')
	if f.BuildV {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = buf
		cmd.Stderr = buf
	}

	if f.BuildWork && tmpdir != "" {
		if runtime.GOOS == "windows" {
			cmd.Env = append(cmd.Env, `TEMP=`+tmpdir)
			cmd.Env = append(cmd.Env, `TMP=`+tmpdir)
		} else {
			cmd.Env = append(cmd.Env, `TMPDIR=`+tmpdir)
		}
	}

	if f.ShouldRun() {
		cmd.Env = Environ(cmd.Env)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s failed: %v%s", strings.Join(cmd.Args, " "), err, buf)
		}
	}
	return nil
}

// environ merges os.Environ and the given "key=value" pairs.
// If a key is in both os.Environ and kv, kv takes precedence.
func Environ(kv []string) []string {
	cur := os.Environ()
	new := make([]string, 0, len(cur)+len(kv))
	goos := runtime.GOOS

	envs := make(map[string]string, len(cur))
	for _, ev := range cur {
		elem := strings.SplitN(ev, "=", 2)
		if len(elem) != 2 || elem[0] == "" {
			// pass the env var of unusual form untouched.
			// e.g. Windows may have env var names starting with "=".
			new = append(new, ev)
			continue
		}
		if goos == "windows" {
			elem[0] = strings.ToUpper(elem[0])
		}
		envs[elem[0]] = elem[1]
	}
	for _, ev := range kv {
		elem := strings.SplitN(ev, "=", 2)
		if len(elem) != 2 || elem[0] == "" {
			panic(fmt.Sprintf("malformed env var %q from input", ev))
		}
		if goos == "windows" {
			elem[0] = strings.ToUpper(elem[0])
		}
		envs[elem[0]] = elem[1]
	}
	for k, v := range envs {
		new = append(new, k+"="+v)
	}
	return new
}

// Creates a new temporary directory. Don't forget to remove.
func NewTmpDir(f *Flags, path string) (string, error) {
	// Make $GOPATH/pkg/work
	tmpdir := ""
	if f.ShouldRun() {
		var err error
		tmpdir, err = ioutil.TempDir(path, "gomobile-work-")
		if err != nil {
			return "", err
		}
	} else {
		if path == "" {
			tmpdir = "$WORK"
		} else {
			tmpdir = filepath.Join(path, "work")
		}
	}
	if f.ShouldPrint() || f.BuildWork {
		fmt.Fprintln(os.Stderr, "WORK="+tmpdir)
	}
	return tmpdir, nil
}

// Returns the directory for a given package.
func PackageDir(f *Flags, pkgpath string) (string, error) {
	pkg, err := build.Default.Import(pkgpath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

func RemoveAll(f *Flags, path string) error {
	if f.ShouldPrint() {
		fmt.Fprintf(os.Stderr, "rm -r -f %s\n", path)
	}
	if f.ShouldRun() {
		return os.RemoveAll(path)
	}
	return nil
}

func WriteFile(flags *Flags, filename string, generate func(io.Writer) error) error {
	if err := Mkdir(flags, filepath.Dir(filename)); err != nil {
		return err
	}
	if flags.ShouldPrint() {
		fmt.Fprintf(os.Stderr, "write %s\n", filename)
	}
	if flags.ShouldRun() {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}()
		return generate(f)
	}
	return generate(ioutil.Discard)
}

func ReadFile(flags *Flags, filename string) ([]byte, error) {
	if flags.ShouldPrint() {
		fmt.Fprintf(os.Stderr, "read %s\n", filename)
	}
	if flags.ShouldRun() {
		return ioutil.ReadFile(filename)
	}
	return []byte{}, nil
}

func CopyFile(f *Flags, dst, src string) error {
	if f.ShouldPrint() {
		fmt.Fprintf(os.Stderr, "cp %s %s\n", src, dst)
	}
	return WriteFile(f, dst, func(w io.Writer) error {
		if f.ShouldRun() {
			f, err := os.Open(src)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(w, f); err != nil {
				return fmt.Errorf("cp %s %s failed: %v", src, dst, err)
			}
		}
		return nil
	})
}

func CopyDir(f *Flags, dst, src string) error {
	cmd := exec.Command("cp", "-R", src, dst)
	return RunCmd(f, "", cmd)
}

func CopyDirContents(f *Flags, dst, src string) error {
	cmd := exec.Command("cp", "-R", src+string(filepath.Separator)+".", dst)
	return RunCmd(f, "", cmd)
}

func Mkdir(flags *Flags, dir string) error {
	if flags.ShouldPrint() {
		fmt.Fprintf(os.Stderr, "mkdir -p %s\n", dir)
	}
	if flags.ShouldRun() {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func Symlink(flags *Flags, src, dst string) error {
	if flags.ShouldPrint() {
		fmt.Fprintf(os.Stderr, "ln -s %s %s\n", src, dst)
	}
	if flags.ShouldRun() {
		// if goos == "windows" {
		//  return doCopyAll(dst, src)
		// }
		return os.Symlink(src, dst)
	}
	return nil
}
