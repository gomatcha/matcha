// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func RunCmd(f *Flags, tmpdir string, cmd *exec.Cmd) error {
	_, err := OutputCmd(f, nil, tmpdir, cmd)
	return err
}

func OutputCmd(f *Flags, fallback []byte, tmpdir string, cmd *exec.Cmd) ([]byte, error) {
	if f.ShouldPrint() {
		str := ""
		if cmd.Dir != "" {
			str += "PWD=" + cmd.Dir + " "
		}
		if len(cmd.Env) > 0 {
			str += strings.Join(cmd.Env, " ") + " "
		}
		str += strings.Join(cmd.Args, " ")
		f.Logger.Println(str)
	}

	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	if f.BuildWork && tmpdir != "" {
		if runtime.GOOS == "windows" {
			cmd.Env = append(cmd.Env, `TEMP=`+tmpdir)
			cmd.Env = append(cmd.Env, `TMP=`+tmpdir)
		} else {
			cmd.Env = append(cmd.Env, `TMPDIR=`+tmpdir)
		}
	}

	var output []byte
	if f.ShouldRun() {
		cmd.Env = MergeEnviron(cmd.Env, os.Environ())
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("%s failed: %v\n%s\n%s", strings.Join(cmd.Args, " "), err, outbuf, errbuf)
		}
		output = outbuf.Bytes()
	} else {
		output = fallback
	}

	if f.BuildV {
		// f.Logger.Println(outbuf.Bytes())
		// f.Logger.Println(errbuf.Bytes())
		if _, err := outbuf.WriteTo(os.Stderr); err != nil {
			return nil, err
		}
		if _, err := outbuf.WriteTo(os.Stdout); err != nil {
			return nil, err
		}
	}
	return output, nil
}

// environ merges os.Environ and the given "key=value" pairs.
// If a key is in both curr and kv, kv takes precedence.
func MergeEnviron(kv, cur []string) []string {
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
		if goos == "windows" { // Windows is case-insensitive?
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
		f.Logger.Println("WORK=" + tmpdir)
	}
	return tmpdir, nil
}

// Returns the directory for a given package.
func PackageDir(f *Flags, pkgpath string) (string, error) {
	if f.ShouldPrint() {
		f.Logger.Printf("go findpackage %s\n", pkgpath)
	}
	if f.ShouldRun() {
		pkg, err := build.Default.Import(pkgpath, "", build.FindOnly)
		if err != nil {
			return "", err
		}
		return pkg.Dir, nil
	}
	return "$GOPATH/src/" + pkgpath, nil
}

func RemoveAll(f *Flags, path string) error {
	if f.ShouldPrint() {
		f.Logger.Printf("rm -r -f %s\n", path)
	}
	if f.ShouldRun() {
		return os.RemoveAll(path)
	}
	return nil
}

func WriteFile(f *Flags, filename string, r io.Reader) (err error) {
	if f.ShouldPrint() {
		f.Logger.Printf("write %s\n", filename)
	}

	disablePrint := f.disablePrint
	f.disablePrint = true
	defer func() {
		f.disablePrint = disablePrint
	}()

	if err = Mkdir(f, filepath.Dir(filename)); err != nil {
		return
	}
	if f.ShouldRun() {
		var file *os.File
		file, err = os.Create(filename)
		if err != nil {
			return
		}
		defer func() {
			if cerr := file.Close(); err == nil {
				err = cerr
			}
		}()

		if _, err = io.Copy(file, r); err != nil {
			return
		}
	}
	return
}

func ReadFile(f *Flags, filename string) ([]byte, error) {
	if f.ShouldPrint() {
		f.Logger.Printf("read %s\n", filename)
	}
	if f.ShouldRun() {
		return ioutil.ReadFile(filename)
	}
	return []byte{}, nil
}

func CopyFile(f *Flags, dst, src string) error {
	if f.ShouldPrint() {
		f.Logger.Printf("cp %s %s\n", src, dst)
	}

	disablePrint := f.disablePrint
	f.disablePrint = true
	defer func() {
		f.disablePrint = disablePrint
	}()

	if f.ShouldRun() {
		file, err := os.Open(src)
		if err != nil {
			return err
		}
		defer file.Close()
		return WriteFile(f, dst, file)
	}
	return nil
}

// func CopyDir(f *Flags, dst, src string) error {
// 	cmd := exec.Command("cp", "-R", src, dst)
// 	return RunCmd(f, "", cmd)
// }

// func CopyDirContents(f *Flags, dst, src string) error {
// 	cmd := exec.Command("cp", "-R", src+string(filepath.Separator)+".", dst)
// 	return RunCmd(f, "", cmd)
// }

func Mkdir(f *Flags, dir string) error {
	if f.ShouldPrint() {
		f.Logger.Printf("mkdir -p %s\n", dir)
	}
	if f.ShouldRun() {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func LookPath(f *Flags, file string) (string, error) {
	if f.ShouldPrint() {
		f.Logger.Printf("which %s\n", file)
	}
	if f.ShouldRun() {
		return exec.LookPath(file)
	}
	return file, nil
}

func GetEnv(f *Flags, key string) string {
	if f.ShouldPrint() {
		f.Logger.Printf("printenv %s\n", key)
	}
	if f.ShouldRun() {
		return os.Getenv(key)
	}
	return "$" + key
}

func ReadDirNames(f *Flags, path string) ([]string, error) {
	if f.ShouldPrint() {
		f.Logger.Printf("ls %s\n", path)
	}
	if f.ShouldRun() {
		file, err := os.Open(path)
		if err != nil {
			return []string{}, err
		}
		defer file.Close()

		return file.Readdirnames(-1)
	}
	return []string{}, nil
}

func IsFile(f *Flags, path string) bool {
	if f.ShouldPrint() {
		f.Logger.Printf("test -f %s\n", path)
	}
	if f.ShouldRun() {
		if st, err := os.Stat(path); err != nil || st.IsDir() {
			return false
		}
	}
	return true
}

func IsDir(f *Flags, path string) bool {
	if f.ShouldPrint() {
		f.Logger.Printf("test -d %s\n", path)
	}
	if f.ShouldRun() {
		if st, err := os.Stat(path); err != nil || !st.IsDir() {
			return false
		}
	}
	return true
}

func Getwd(f *Flags) (string, error) {
	if f.ShouldPrint() {
		f.Logger.Printf("pwd\n")
	}
	if f.ShouldRun() {
		return os.Getwd()
	}
	return "$CWD", nil
}
