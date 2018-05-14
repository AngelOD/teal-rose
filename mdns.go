package main

import (
	"net"
	"time"

	"github.com/hashicorp/mdns"
)

func setupMdns() {
	var service *mdns.MDNSService
	var err error

	info := []string{"Test service"}
	if len(webIP) > 0 {
		var ips []net.IP

		ips = append(ips, net.ParseIP(webIP))
		service, err = mdns.NewMDNSService(domain, mdnsServiceType, webDomain, webHost, webPort, ips, info)
	} else {
		service, err = mdns.NewMDNSService(domain, mdnsServiceType, webDomain, webHost, webPort, nil, info)
		logger.Infof("Inside: %+v", service)
	}

	if err != nil {
		logger.Errorf("Unable to create mDNS service: %v", err)
		return
	}

	logger.Infof("Outside: %+v", service)
	logger.Infof("IPs: %+v", service.IPs)

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		logger.Errorf("Unable to start mDNS server: %v", err)
		return
	}
	defer server.Shutdown()

	for {
		select {
		case <-prg.exit:
			logger.Info("Shutting down error logging interface.")
			return
		default:
			time.Sleep(5 * time.Second)
		}
	}
}
