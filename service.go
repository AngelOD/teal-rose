package main

import (
	"github.com/kardianos/service"
	"log"
	"os"
)

type Config struct {
	Name, DisplayName, Description string
}

var logger service.Logger
var svc service.Service
var prg *program

type program struct {
	exit    chan struct{}
	service service.Service
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager")
	}

	p.exit = make(chan struct{})
	go p.run()

	return nil
}

func (p *program) run() {
	logger.Infof("I'm running. [%v]", service.Platform())

	if debugLog && !service.Interactive() {
		fileName, err := getConfigPath("output.log")
		if err == nil {
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				logger.Infof("Using log file: %s", fileName)
				defer f.Close()
				log.SetOutput(f)
			} else {
				logger.Warning(err)
			}
		} else {
			logger.Warning(err)
		}
	} else if debugLog {
		logger.Info("Outputting log to standard out.")
	}

	SocketServer(port, p)
}

func (p *program) Stop(s service.Service) error {
	logger.Info("I'm stopping!")
	close(p.exit)
	return nil
}

func initService() {
	svcConfig := &service.Config{
		Name:        "Sw802f18Receiver",
		DisplayName: "SW802F18 Sensor Data Receiver",
		Description: "Listens on a specific port and receives data from sensor clusters.",
		Arguments:   []string{"run", "-s", "-d"},
	}

	prg = &program{}
	s, err := service.New(prg, svcConfig)

	if err != nil {
		log.Fatalln(err)
	}

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)

	if err != nil {
		log.Fatalln(err)
	}

	prg.service = s

	go func() {
		for {
			select {
			case err := <-errs:
				if err != nil {
					logger.Error(err)
				}
			case <-prg.exit:
				logger.Info("Shutting down error logging interface.")
				return
			default:
				// Do nothing
			}
		}
	}()

	svc = s
}
