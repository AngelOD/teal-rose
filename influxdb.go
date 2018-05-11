package main

import (
	"log"
	"time"
	"github.com/influxdata/influxdb/client/v2"
)

func runInfluxDbTest() {
	c, err := client.NewHTTPClient(client.HTTPConfig {
		Addr: "http://localhost:8086",
		Username: "",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig {
		Database: "",
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle": 10.1,
		"system": 53.3,
		"user": 46.6,
	}

	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}
	bp.AddPoint(pt)

	if err := c.Write(bp); err != nil {
		log.Fatal(err)
		return
	}

	if err := c.Close(); err != nil {
		log.Fatal(err)
		return
	}
}