package main

import (
	"fmt"
	"github.com/teris-io/cli"
)

func setupCli() cli.App {
	optPort := cli.NewOption("port", "Port to host on").WithType(cli.TypeInt)
	optStoreData := cli.NewOption("save", "Store information in DB").WithChar('s').WithType(cli.TypeBool)

	app := cli.New("SW802F18 Test Server").WithOption(optPort).WithOption(optStoreData).WithAction(handleCli)

	return app
}

func handleCli(args []string, options map[string]string) int {
	fmt.Println("Args:", args)
	fmt.Println("Options:", options)

	TestJson()

	return 0
}
