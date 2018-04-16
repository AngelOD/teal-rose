// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// for the current OS and architecture
func Build() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild(runtime.GOOS, runtime.GOARCH, false)
}

// for debugging on the current OS and architecture
func BuildDebug() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild(runtime.GOOS, runtime.GOARCH, true)
}

// for all supported operating systems and architectures
func BuildAll() {
	mg.Deps(BuildAll32)
	mg.Deps(BuildAll64)
}

// for all supported operating systems in 64 bit versions
func BuildAll64() {
	mg.Deps(BuildLinux)
	mg.Deps(BuildWindows)
}

// for all supported operating systems in 32 bit versions
func BuildAll32() {
	mg.Deps(BuildLinux32)
	mg.Deps(BuildWindows32)
}

// for the amd64 architecture on Linux
func BuildLinux() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild("linux", "amd64", false)
}

// for the 386 architecture on Linux
func BuildLinux32() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild("linux", "386", false)
}

// for the amd64 architecture on Windows
func BuildWindows() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild("windows", "amd64", false)
}

// for the 386 architecture on Windows
func BuildWindows32() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild("windows", "386", false)
}

// runs all tests
func Test() {
	fmt.Println("Running tests...")

	sh.RunV("go", "test", "-v")
}

func doBuild(pOs string, pArch string, pDebugBuild bool) {
	fmt.Printf("Building for %s (%s)\n", pOs, pArch)

	path := filepath.Join(getOutputPath(), fmt.Sprintf("%s_%s", pOs, pArch))
	if strings.ToLower(pOs) == "windows" {
		path += ".exe"
	}

	fmt.Printf("Outputting to: %s\n", path)

	os.Setenv("GOOS", pOs)
	os.Setenv("GOARCH", pArch)

	debugFlags := ""
	if !pDebugBuild {
		fmt.Println("Stripping debug info from executable")
		debugFlags = "-ldflags=-s -w"
	}

	sh.RunV("go", "build", "-o", path, debugFlags)
}

func doGenerate() {
	sh.RunV("go", "generate")
}

func getOutputPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error! %v", err))
	}

	path := filepath.Join(dir, "_dist")

	if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
		panic(fmt.Sprintf("Error! %v", err))
	}

	return path
}
