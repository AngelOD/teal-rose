package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func storeDataRunner() {
	var rd radioData
	var sd sensorDataCombined

	logger.Info("Starting data storage thread.")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbData["DB_USER"], dbData["DB_PASS"], dbData["DB_NAME"]))

	if err != nil {
		logger.Errorf("Error connecting to DB: %s", err)
		// TODO: Figure out a good way to exit here
		return
	}
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO radio_datas (" +
		"`radio_bus_id`, `channel`, `node_mac_address`, `packet_type`, `sequence_number`," +
		"`timestamp_nano`, `timestamp_tz`, `v_bat`, `vcc`, `temperature`," +
		"`humidity`, `pressure`, `co2`, `tvoc`, `light`," +
		"`uv`, `sound_pressure`, `port_input`, `mag`, `acc`," +
		"`gyro`) VALUES (" +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmtInsert.Close()

	if err != nil {
		logger.Errorf("Error preparing statement: %s", err)
		// TODO: Figure out a good way to exit here
		return
	}

	logger.Infof("Statement prepared. Listening to save requests...\n")

	for {
		select {
		case rd = <-rdStore:
			logger.Info("Data found!")

			sd = rd.GetSensorData()

			if shouldDiscardData(&sd) {
				logger.Warning("Invalid data. Ignoring row.")
				continue
			}

			t, err := time.Parse(rfc3339Micro, rd.TimestampTz)
			if err != nil {
				logger.Warningf("Unable to convert time string: %s", rd.TimestampTz)
				logger.Warningf("Because of: %s", err)
				continue
			}

			_, err = stmtInsert.Exec(
				rd.RadioBusID, rd.Channel, rd.NodeMacAddress, rd.PacketType, rd.SequenceNumber,
				t.UnixNano(), rd.TimestampTz, sd.vBat, sd.vcc, sd.temperature,
				sd.humidity, sd.pressure, sd.co2, sd.tvoc, sd.light,
				sd.uv, sd.soundPressure, sd.portInput, sd.mag.String(), sd.acc.String(),
				sd.gyro.String(),
			)

			if err != nil {
				logger.Errorf("ERROR! Unable to save data: %s", err)
				continue
			}
		case <-prg.exit:
			logger.Info("Shutting down storeDataRunner.")
			return
		default:
			time.Sleep(10 * time.Second)
		}
	}
}

func ensureDataTable() bool {
	// TODO Complete this
	return true
}

func shouldDiscardData(sd *sensorDataCombined) bool {
	return sd.pressure == 0 || sd.humidity == 0 || sd.temperature == 0 || sd.co2 == 0
}
