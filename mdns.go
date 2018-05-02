package main

import (
	"net"
	"time"

	"github.com/hashicorp/mdns"
)

func setupMdns() {
	var ips []net.IP

	if len(webIP) > 0 {
		ips = append(ips, net.ParseIP(webIP))
	} else {
		ips = nil
	}

	info := []string{"Test service"}
	service, err := mdns.NewMDNSService(domain, mdnsServiceType, webDomain, webHost, webPort, ips, info)

	logger.Infof("IPs: %+v", service.IPs)

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
