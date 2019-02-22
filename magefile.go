// Copyright 2019 Wio. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build mage

package main

import (
	"errors"
	"fmt"
	"github.com/magefile/mage/sh"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	packageName = "github.com/wio/wio"
	execName    = "wio"
)

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	os.Setenv("GO111MODULE", "on")
}

// Pull required tools
func Setup() error {
	os.Setenv("GO111MODULE", "off")

	fmt.Println("Getting all the go tools needed...")

	if runtime.GOOS == "windows" {
		if err := sh.RunWith(nil, goexe, "get", "-u", "golang.org/x/sys/windows"); err != nil {
			return err
		}
	} else {
		if err := sh.RunWith(nil, goexe, "get", "-u", "golang.org/x/sys/unix"); err != nil {
			return err
		}
	}

	// when TRAVIS is not running
	if os.Getenv("TRAVIS") == "" {
		if err := sh.RunWith(nil, goexe, "get", "-u", "golang.org/x/tools/..."); err != nil {
			return err
		}
	}

	if err := sh.RunWith(nil, goexe,  "get", "-u", "github.com/davecgh/go-spew/spew"); err != nil {
		return err
	}
	if err := sh.RunWith(nil, goexe, "get", "-u", "github.com/stretchr/testify"); err != nil {
		return err
	}
	if err := sh.RunWith(nil, goexe, "get", "-u", "github.com/pmezard/go-difflib/difflib"); err != nil {
		return err
	}
	if err := sh.RunWith(nil, goexe, "get", "-u", "github.com/kevinburke/go-bindata/..."); err != nil {
		return err
	}

	return nil
}

// Build wio binary
func Build() error {
	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	os.Setenv("GO111MODULE", "on")

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(currDir + "/pkg/util/sys"); err != nil {
		return err
	}

	if err := sh.RunWith(nil, "go-bindata", "-nomemcopy", "-pkg", "sys",
		"-prefix", "../../../", "../../../assets/..."); err != nil {
		return err
	}

	if err := os.Chdir(currDir + "/cmd/" + execName); err != nil {
		return err
	}

	if err := sh.RunWith(nil, goexe, "build", "-ldflags=-s -w", "-o",
		"../../bin/" + execName, "-v"); err != nil {
		return err
	}

	return nil
}

// Clean binary files
func Clean() error {
	return sh.RunWith(nil, goexe, "clean")
}

func Test() error {
	if runtime.GOOS == "windows" {
		return errors.New("running tests not supported on windows")
	}

	err := os.Chdir("test")
	if err != nil {
		return err
	}

	if err := sh.RunWith(nil, "bash", "./symlinks.sh"); err != nil {
		return err
	}
	return sh.RunWith(nil, "bash", "./runtests.sh")
}

// Format code
func Fmt() error {
	return sh.RunWith(nil, goexe, "fmt", "./...")
}

// Install dependencies
func Install() error {
	return sh.RunWith(nil, goexe, "install")
}

var (
	pkgPrefixLen = len(execName)
	pkgs         []string
	pkgsInit     sync.Once
)

// Run gofmt linter
func FmtLint() error {
	if !isGoLatest() {
		return errors.New("Go version must be 1.11")
	}
	pkgs, err := wioPackages()
	if err != nil {
		return err
	}
	failed := false
	first := true
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			// gofmt doesn't exit with non-zero when it finds unformatted code
			// so we have to explicitly look for output, and if we find any, we
			// should fail this target.
			s, err := sh.Output("gofmt", "-l", f)
			if err != nil {
				fmt.Printf("ERROR: running gofmt on %q: %v\n", f, err)
				failed = true
			}
			if s != "" {
				if first {
					fmt.Println("The following files are not gofmt'ed:")
					first = false
				}
				failed = true
				fmt.Println(s)
			}
		}
	}
	if failed {
		return errors.New("improperly formatted go files")
	}
	return nil
}

func wioPackages() ([]string, error) {
	var err error
	pkgsInit.Do(func() {
		var s string
		s, err = sh.Output(goexe, "list", "./...")
		if err != nil {
			return
		}
		pkgs = strings.Split(s, "\n")
		for i := range pkgs {
			pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
		}
	})
	return pkgs, err
}

func isGoLatest() bool {
	return strings.Contains(runtime.Version(), "1.11")
}
