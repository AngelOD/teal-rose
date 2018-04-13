package main

import "github.com/hashicorp/mdns"

func setupMdns() {
	info := []string{"Test service"}
	service, err := mdns.NewMDNSService(MdnsServiceType + "_" + domain, MdnsServiceType, domain, host, port, nil, info)

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
}
