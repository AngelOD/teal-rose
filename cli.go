package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alexsasharegan/dotenv"
	"github.com/blang/semver"
	"github.com/kardianos/osext"
	"github.com/kardianos/service"
	"github.com/teris-io/cli"
)

func setupCli() cli.App {
	optDebug := cli.NewOption("debug", "Log debug information").WithChar('d').WithType(cli.TypeBool)
	optDomain := cli.NewOption("domain", "Domain name of service").WithChar('D').WithType(cli.TypeString)
	optHost := cli.NewOption("host", "Host name of current host").WithChar('H').WithType(cli.TypeString)
	optIDbHost := cli.NewOption("idbhost", "Influx DB host").WithType(cli.TypeString)
	optIDbName := cli.NewOption("idbname", "Influx DB name").WithType(cli.TypeString)
	optIDbPass := cli.NewOption("idbpass", "Influx DB pass").WithType(cli.TypeString)
	optIDbUser := cli.NewOption("idbuser", "Influx DB user").WithType(cli.TypeString)
	optMysqlConn := cli.NewOption("myconn", "MySQL connector string").WithType(cli.TypeString)
	optPort := cli.NewOption("port", "Port to host on").WithType(cli.TypeInt)
	optStoreData := cli.NewOption("save", "Store information in DB").WithChar('s').WithType(cli.TypeBool)

	cmdFix := cli.NewCommand("fix", "Clean data of dual results").
		WithAction(handleFixCli)

	cmdRun := cli.NewCommand("run", "Normal run").
		WithOption(optDebug).
		WithOption(optDomain).
		WithOption(optHost).
		WithOption(optPort).
		WithOption(optStoreData).
		WithAction(handleCli)

	cmdService := cli.NewCommand("service", "Manage service").
		WithShortcut("srv").
		WithArg(cli.NewArg("command", "The service subcommand")).
		WithAction(handleServerCli)

	cmdTestInflux := cli.NewCommand("tin", "Migrate to InfluxDB").
		WithOption(optIDbHost).
		WithOption(optIDbName).
		WithOption(optIDbUser).
		WithOption(optIDbPass).
		WithOption(optMysqlConn).
		WithAction(handleTinCli)

	cmdVersion := cli.NewCommand("version", "Version info").
		WithShortcut("ver").
		WithAction(handleVersionCli)

	app := cli.New("SW802F18 Test Server").
		WithCommand(cmdFix).
		WithCommand(cmdService).
		WithCommand(cmdRun).
		WithCommand(cmdVersion).
		WithCommand(cmdTestInflux)

	return app
}

func parseCommonCli(args []string, options map[string]string) int {
	if pDebug, err := strconv.ParseBool(options["debug"]); err == nil && pDebug {
		debugLog = true
	}

	if pDomain, prs := options["domain"]; prs {
		domain = pDomain
	} else {
		domain, _ = os.Hostname()
	}

	if pHost, prs := options["host"]; prs {
		host = pHost
	}

	if pPort, err := strconv.Atoi(options["port"]); err == nil && pPort > 1024 && pPort <= 65535 {
		port = pPort
	}

	return 0
}

func parseDbCli(args []string, options map[string]string) int {
	keys := []string{"DB_TYPE", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASS"}
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

	return 0
}

func handleCli(args []string, options map[string]string) int {
	if !loadDotEnv() {
		return 2
	}

	if common := parseCommonCli(args, options); common > 0 {
		return common
	}

	if pSave, err := strconv.ParseBool(options["save"]); err == nil && pSave {
		if dbRes := parseDbCli(args, options); dbRes > 0 {
			return dbRes
		}
	}

	logger.Info("Using port number: ", port)

	if debugLog {
		logger.Info("Debug logging enabled")
	} else {
		logger.Info("Debug logging disabled")
	}

	if saveData {
		if strings.ToLower(dbData["DB_TYPE"]) == "influxdb" {
			go influxStoreDataRunner()
		} else {
			go mysqlStoreDataRunner()
		}

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

func handleFixCli(args []string, options map[string]string) int {
	if !loadDotEnv() {
		return 2
	}

	if common := parseCommonCli(args, options); common > 0 {
		return common
	}

	if dbRes := parseDbCli(args, options); dbRes > 0 {
		return dbRes
	}

	if err := fixDatabaseEntries(); err != nil {
		return 2
	}

	return 0
}

func handleTinCli(args []string, options map[string]string) int {
	var (
		influxDbHost   string
		influxDbName   string
		influxDbUser   string
		influxDbPass   string
		mysqlConnector string
		prs            bool
	)

	if influxDbHost, prs = options["idbhost"]; !prs {
		logger.Error("Missing option: Influx DB host")
		return 2
	}

	if influxDbName, prs = options["idbname"]; !prs {
		logger.Error("Missing option: Influx DB name")
		return 2
	}

	if influxDbUser, prs = options["idbuser"]; !prs {
		logger.Error("Missing option: Influx DB user")
		return 2
	}

	if influxDbPass, prs = options["idbpass"]; !prs {
		logger.Error("Missing option: Influx DB pass")
		return 2
	}

	if mysqlConnector, prs = options["myconn"]; !prs {
		logger.Error("Missing option: MySql DB connector string")
		return 2
	}

	migrateMysqlToInflux(mysqlConnector, influxDbHost, influxDbName, influxDbUser, influxDbPass)

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

func handleVersionCli(args []string, options map[string]string) int {
	vers, err := semver.ParseTolerant(versionInfo)
	if err != nil {
		logger.Errorf("Unable to parse version info: %v", err)
		return 1
	}

	logger.Infof("Version: v%d.%d.%d", vers.Major, vers.Minor, vers.Patch)

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

	// Handle the vars that needs to be possible to apply through .env
	if debugVal, prs := env["DEBUG"]; prs {
		if pDebugVal, err := strconv.ParseBool(debugVal); err == nil {
			debugLog = pDebugVal
		}
	}

	if webDomainVal, prs := env["WEB_DOMAIN"]; prs {
		webDomain = webDomainVal
	}

	if webHostVal, prs := env["WEB_HOST"]; prs {
		webHost = webHostVal
	}

	if webIPVal, prs := env["WEB_IP"]; prs {
		webIP = webIPVal
	}

	if webPortVal, prs := env["WEB_PORT"]; prs {
		if pWebPortVal, err := strconv.Atoi(webPortVal); err == nil {
			webPort = pWebPortVal
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
