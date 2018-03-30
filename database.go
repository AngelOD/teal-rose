package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func StoreData(rds []RadioData) bool {
	if len(rds) < 1 {
		return true
	}

	var rd RadioData
	var sd SensorDataCombined
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbData["DB_USER"], dbData["DB_PASS"], dbData["DB_NAME"]))

	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err)
		return false
	}
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO radio_datas (" +
		"`radio_bus_id`, `channel`, `node_mac_address`, `packet_type`, `sequence_number`," +
		"`timestamp`, `timestamp_tz`, `v_bat`, `vcc`, `temperature`," +
		"`humidity`, `pressure`, `co2`, `tvoc`, `light`," +
		"`uv`, `sound_pressure`, `port_input`, `mag`, `acc`," +
		"`gyro`) VALUES (" +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmtInsert.Close()

	if err != nil {
		log.Fatalf("Error preparing statement: %s", err)
		return false
	}

	log.Printf("Statement prepared. Inserting data...\n")

	for i := 0; i < len(rds); i++ {
		rd = rds[i]
		sd = rd.GetSensorData()
		_, err = stmtInsert.Exec(
			rd.RadioBusId, rd.Channel, rd.NodeMacAddress, rd.PacketType, rd.SequenceNumber,
			rd.Timestamp, rd.TimestampTz, sd.VBat, sd.Vcc, sd.Temperature,
			sd.Humidity, sd.Pressure, sd.Co2, sd.Tvoc, sd.Light,
			sd.Uv, sd.SoundPressure, sd.PortInput, sd.Mag.String(), sd.Acc.String(),
			sd.Gyro.String(),
		)

		if err != nil {
			log.Fatalf("ERROR! Failed after %d entries: %s", i, err)
			return false
		}
	}

	return true
}

func EnsureDataTable() bool {
	// TODO Complete this
	return true
}
