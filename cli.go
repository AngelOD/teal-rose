package main

import (
	"path/filepath"
	"strconv"

	"github.com/alexsasharegan/dotenv"
	"github.com/kardianos/osext"
	"github.com/kardianos/service"
	"github.com/teris-io/cli"
)

var (
	port      = 3333
	debugLog  = false
	saveData  = false
	saveEvery = 5
	dbData    = map[string]string{
		"DB_NAME": "teal_rose",
		"DB_USER": "teal_rose",
		"DB_PASS": "",
	}
)

var env map[string]string;

func setupCli() cli.App {
	optDebug := cli.NewOption("debug", "Log debug information").WithChar('d').WithType(cli.TypeBool)
	optPort := cli.NewOption("port", "Port to host on").WithType(cli.TypeInt)
	optStoreData := cli.NewOption("save", "Store information in DB").WithChar('s').WithType(cli.TypeBool)

	cmdService := cli.NewCommand("service", "Manage service").
		WithShortcut("srv").
		WithArg(cli.NewArg("command", "The service subcommand")).
		WithAction(handleServerCli)

	cmdRun := cli.NewCommand("run", "Normal run").
		WithOption(optDebug).
		WithOption(optPort).
		WithOption(optStoreData).
		WithAction(handleCli)

	app := cli.New("SW802F18 Test Server").
		WithCommand(cmdService).
		WithCommand(cmdRun)

	return app
}

func handleCommonCli(args []string, options map[string]string) int {
	if pPort, err := strconv.Atoi(options["port"]); err == nil && pPort > 1024 && pPort <= 65535 {
		port = pPort
	}

	if pDebug, err := strconv.ParseBool(options["debug"]);err == nil && pDebug {
		debugLog = true
	}

	return 0
}

func handleCli(args []string, options map[string]string) int {
	if !loadDotEnv() {
		return 2
	}

	common := handleCommonCli(args, options)

	if common > 0 {
		return common
	}

	if pSave, err := strconv.ParseBool(options["save"]); err == nil && pSave {
		keys := []string{"DB_NAME", "DB_USER", "DB_PASS"}
		saveData = true

		for i := 0; i < len(keys); i++ {
			if val, prs := env[keys[i]]; prs && len(val) > 0 {
				dbData[keys[i]] = val
			} else {
				if len(dbData[keys[i]]) == 0 {
					logger.Errorf("Variable %s required but missing from .env file!", keys[i])
					return 2
				}
			}
		}
	}

	logger.Info("Using port number: ", port)

	if debugLog {
		logger.Info("Debug logging enabled")
	} else {
		logger.Info("Debug logging disabled")
	}

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
		return 2
	}

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

func loadDotEnv() bool {
	envPath, err := getConfigPath(".env")
	if err != nil {
		logger.Error(err)
		return false
	}

	env, err = dotenv.ReadFile(envPath)
	if err != nil {
		logger.Error(err)
		return false
	}

	if debugVal, prs := env["DEBUG"]; prs {
		if pDebugVal, err := strconv.ParseBool(debugVal); err == nil {
			debugLog = pDebugVal
		}
	}

	return true
}

func getConfigPath(fileName string) (string, error) {
	fullexecpath, err := osext.Executable()
	if err != nil {
		return "", err
	}

	dir, _ := filepath.Split(fullexecpath)

	return filepath.Join(dir, fileName), nil
}
