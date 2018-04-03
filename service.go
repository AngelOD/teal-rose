package main

import (
	"github.com/kardianos/service"
	"time"
	"log"
)

type Config struct {
	Name, DisplayName, Description string

	Stderr, Stdout string
}

var logger service.Logger

type program struct {
	exit chan struct{}
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

func (p *program) run() error {
	logger.Infof("I'm running. [%v]", service.Platform())
	ticker := time.NewTicker(2 * time.Second)

	for {
		logger.Info("Checking...")

		select {
		case tm := <-ticker.C:
			logger.Infof("Still running at %v...", tm)
		case <-p.exit:
			logger.Info("Triggered by exit!")
			ticker.Stop()
			return nil
		}
	}
}

func (p *program) Stop(s service.Service) error {
	logger.Info("I'm stopping!")
	close(p.exit)
	return nil
}

func initService() service.Service {
	svcConfig := &service.Config{
		Name: "Sw802f18Receiver",
		DisplayName: "SW802F18 Sensor Data Receiver",
		Description: "Listens on a specific port and receives data from sensor clusters.",
		Arguments: []string{"-s"},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)

	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)

	if err != nil {
		log.Fatal(err)
	}

	prg.service = s

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	return s
}