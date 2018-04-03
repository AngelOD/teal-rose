package main

import (
	"os"
	"strconv"

	"github.com/alexsasharegan/dotenv"
	"github.com/kardianos/service"
	"github.com/teris-io/cli"
)

var (
	port      = 3333
	saveData  = false
	saveEvery = 10
	dbData    = map[string]string{
		"DB_NAME": "teal_rose",
		"DB_USER": "teal_rose",
		"DB_PASS": "",
	}
)

func setupCli() cli.App {
	cmdService := cli.NewCommand("service", "Manage service").
		WithShortcut("srv").
		WithArg(cli.NewArg("command", "The service subcommand")).
		WithAction(handleServerCli)

	optPort := cli.NewOption("port", "Port to host on").WithType(cli.TypeInt)
	optStoreData := cli.NewOption("save", "Store information in DB").WithChar('s').WithType(cli.TypeBool)

	app := cli.New("SW802F18 Test Server").
		WithCommand(cmdService).
		WithOption(optPort).
		WithOption(optStoreData).
		WithAction(handleCli)

	return app
}

func handleCli(args []string, options map[string]string) int {
	if pPort, err := strconv.Atoi(options["port"]); err == nil && pPort > 1024 && pPort <= 65535 {
		port = pPort
	}

	if pSave, err := strconv.ParseBool(options["save"]); err == nil && pSave {
		keys := []string{"DB_NAME", "DB_USER", "DB_PASS"}
		saveData = true

		dotenv.Load()

		for i := 0; i < len(keys); i++ {
			if val, ok := os.LookupEnv(keys[i]); ok && len(val) > 0 {
				dbData[keys[i]] = val
			} else {
				if len(dbData[keys[i]]) == 0 {
					logger.Errorf("Variable %s required but missing!", keys[i])
					return 2
				}
			}
		}
	}

	logger.Info("Using port number: ", port)

	if saveData {
		logger.Info("Will save incoming data to DB.")
		logger.Infof("Name: %s\n", dbData["DB_NAME"])
		logger.Infof("User: %s\n", dbData["DB_USER"])
		logger.Infof("Pass: HIDDEN [%d]", len(dbData["DB_PASS"]))
	} else {
		logger.Info("Will discard incoming data.")
	}

	err := svc.Run()

	if err != nil {
		logger.Error(err)
		return 3
	}

	//SocketServer(port)

	return 0
}

func handleServerCli(args []string, options map[string]string) int {
	err := service.Control(svc, args[0])

	if err != nil {
		logger.Infof("Valid actions: %q\n", service.ControlAction)
		logger.Error(err)
	}

	return 0
}
