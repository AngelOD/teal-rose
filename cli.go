package main

import (
	"log"
	"os"
	"strconv"

	"github.com/alexsasharegan/dotenv"
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
	optPort := cli.NewOption("port", "Port to host on").WithType(cli.TypeInt)
	optStoreData := cli.NewOption("save", "Store information in DB").WithChar('s').WithType(cli.TypeBool)

	app := cli.New("SW802F18 Test Server").WithOption(optPort).WithOption(optStoreData).WithAction(handleCli)

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
					log.Fatalf("Variable %s required but missing!", keys[i])
					return 2
				}
			}
		}
	}

	log.Println("Using port number: ", port)

	if saveData {
		log.Println("Will save incoming data to DB.")
		log.Printf("Name: %s\n", dbData["DB_NAME"])
		log.Printf("User: %s\n", dbData["DB_USER"])
		log.Printf("Pass: HIDDEN [%d]", len(dbData["DB_PASS"]))
	} else {
		log.Println("Will discard incoming data.")
	}

	SocketServer(port)

	return 0
}
