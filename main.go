// +build go1.10

package main

import "os"

func main() {
	app := setupCli()

	os.Exit(app.Run(os.Args, os.Stdout))

	//port := 3333
	//SocketServer(port)
}
