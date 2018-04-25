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

var flags string = ""

// for the current OS and architecture
func Build() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild(runtime.GOOS, runtime.GOARCH)
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
	doBuild("linux", "amd64")
}

// for the 386 architecture on Linux
func BuildLinux32() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild("linux", "386")
}

// for the amd64 architecture on Windows
func BuildWindows() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild("windows", "amd64")
}

// for the 386 architecture on Windows
func BuildWindows32() {
	mg.Deps(doGenerate)
	mg.Deps(Test)
	doBuild("windows", "386")
}

// runs all tests
func Test() {
	fmt.Println("Running tests...")

	cPath := filepath.Join(getOutputPath(), "cover.out")
	oPath := filepath.Join(getOutputPath(), "coverage.html")

	sh.RunV("go", "test", "-bench", ".", "-v", "-coverprofile="+cPath)
	sh.RunV("go", "tool", "cover", "-html="+cPath, "-o", oPath)
	sh.RunV("go", "tool", "cover", "-func="+cPath)
	sh.RunV("rm", cPath)

	fmt.Printf("Detailed coverage data available in %s\n", oPath)
}

// runs go clean
func Clean() {
	fmt.Println("Cleaning")

	sh.RunV("go", "clean")
}

func doBuild(pOs string, pArch string) {
	fmt.Printf("Building for %s (%s)\n", pOs, pArch)

	path := filepath.Join(getOutputPath(), fmt.Sprintf("%s_%s", pOs, pArch))
	if strings.ToLower(pOs) == "windows" {
		path += ".exe"
	}

	fmt.Printf("Outputting to: %s\n", path)

	setFlags()

	os.Setenv("GOOS", pOs)
	os.Setenv("GOARCH", pArch)

	sh.RunV("go", "build", "-o", path, "-ldflags="+flags, "github.com/angelod/teal-rose")
}

func doGenerate() {
	sh.RunV("go", "generate")
}

func setFlags() {
	tag := getTag()

	if tag == "" {
		tag = "v0.0.1-dev"
	}

	flags = fmt.Sprintf(`-X "main.versionInfo=%s"`, tag)
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

func getTag() string {
	tag, _ := sh.Output("git", "describe", "--tags")
	return tag
}
