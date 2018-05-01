package main

import (
	"time"

	"github.com/hashicorp/mdns"
)

func setupMdns() {
	info := []string{"Test service"}
	service, err := mdns.NewMDNSService(domain, mdnsServiceType, "", "", port, nil, info)

	if err != nil {
		logger.Errorf("Unable to create mDNS service: %v", err)
		return
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		logger.Errorf("Unable to start mDNS server: %v", err)
		return
	}
	defer server.Shutdown()

	go func() {
		for {
			select {
			case <-prg.exit:
				logger.Info("Shutting down error logging interface.")
				return
			default:
				time.Sleep(5 * time.Second)
			}
		}
	}()
}
