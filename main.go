// +build go1.10

package main

import (
	"os"
)

func main() {
	initService()

	app := setupCli()

	os.Exit(app.Run(os.Args, os.Stdout))
}
