package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"

var rdStore = make(chan RadioData, 20)

func storeDataRunner() {
	var rd RadioData
	var sd SensorDataCombined

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
			sd = rd.GetSensorData()

			if shouldDiscardData(&sd) {
				logger.Warning("Invalid data. Ignoring row.")
				continue
			}

			t, err := time.Parse(RFC3339Micro, rd.TimestampTz)
			if err != nil {
				logger.Warningf("Unable to convert time string: %s", rd.TimestampTz)
				logger.Warningf("Because of: %s", err)
				continue
			}

			_, err = stmtInsert.Exec(
				rd.RadioBusId, rd.Channel, rd.NodeMacAddress, rd.PacketType, rd.SequenceNumber,
				t.UnixNano(), rd.TimestampTz, sd.VBat, sd.Vcc, sd.Temperature,
				sd.Humidity, sd.Pressure, sd.Co2, sd.Tvoc, sd.Light,
				sd.Uv, sd.SoundPressure, sd.PortInput, sd.Mag.String(), sd.Acc.String(),
				sd.Gyro.String(),
			)

			if err != nil {
				logger.Errorf("ERROR! Unable to save data: %s", err)
				continue
			}
		case <-prg.exit:
			return
		default:
			// Do nothing
		}
	}
}

func StoreData(rds []RadioData) bool {
	if len(rds) < 1 {
		return true
	}

	var rd RadioData
	var sd SensorDataCombined
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbData["DB_USER"], dbData["DB_PASS"], dbData["DB_NAME"]))

	if err != nil {
		logger.Errorf("Error connecting to DB: %s", err)
		return false
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
		return false
	}

	logger.Infof("Statement prepared. Inserting data...\n")

	for i := 0; i < len(rds); i++ {
		rd = rds[i]
		sd = rd.GetSensorData()

		if shouldDiscardData(&sd) {
			logger.Warning("Invalid data. Ignoring row.")
			continue
		}

		t, err := time.Parse(RFC3339Micro, rd.TimestampTz)
		if err != nil {
			logger.Warningf("Unable to convert time string: %s", rd.TimestampTz)
			logger.Warningf("Because of: %s", err)
			continue
		}

		_, err = stmtInsert.Exec(
			rd.RadioBusId, rd.Channel, rd.NodeMacAddress, rd.PacketType, rd.SequenceNumber,
			t.UnixNano(), rd.TimestampTz, sd.VBat, sd.Vcc, sd.Temperature,
			sd.Humidity, sd.Pressure, sd.Co2, sd.Tvoc, sd.Light,
			sd.Uv, sd.SoundPressure, sd.PortInput, sd.Mag.String(), sd.Acc.String(),
			sd.Gyro.String(),
		)

		if err != nil {
			logger.Errorf("ERROR! Failed after %d entries: %s", i, err)
			return false
		}
	}

	return true
}

func EnsureDataTable() bool {
	// TODO Complete this
	return true
}

func shouldDiscardData(sd *SensorDataCombined) bool {
	return sd.Pressure == 0 || sd.Humidity == 0 || sd.Temperature == 0 || sd.Co2 == 0
}
