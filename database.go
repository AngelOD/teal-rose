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

	stmtInsert, err := db.Prepare(
		"INSERT INTO radio_datas (" +
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

	stmtSelect, err := db.Prepare(
		"SELECT `co2`, `humidity`, `light`, `pressure`, `sound_pressure`, `temperature`, `tvoc`, `uv` " +
			"FROM radio_datas " +
			"WHERE `node_mac_address` = ? AND `sequence_number` = ?")
	defer stmtSelect.Close()

	if err != nil {
		logger.Errorf("Error preparing statement: %s", err)
		// TODO: Figure out a good way to exit here
		return
	}

	logger.Infof("Statements prepared. Listening for save requests...\n")

	for {
		select {
		case rd = <-rdStore:
			logger.Info("Data found!")

			sd = rd.getSensorData()

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

			rows, err := stmtSelect.Query(rd.NodeMacAddress, rd.SequenceNumber)

			if err != nil {
				logger.Errorf("Error querying data: %s", err)
			} else if rows.Next() {
				var (
					tCo2           int
					tHumidity      int
					tLight         int
					tPressure      int
					tSoundPressure int
					tTemperature   int
					tTvoc          int
					tUv            int
				)

				if err := rows.Scan(&tCo2, &tHumidity, &tLight, &tPressure, &tSoundPressure, &tTemperature, &tTvoc, &tUv); err != nil {
					logger.Errorf("Error parsing data: %s", err)
				} else if tCo2 == sd.co2 && tHumidity == sd.humidity && tLight == sd.light && tPressure == sd.pressure &&
					tSoundPressure == sd.soundPressure && tTemperature == sd.temperature && tTvoc == sd.tvoc && tUv == sd.uv {
					//
					logger.Info("Multiple entries detected. Doing nothing for now.")
				} else {
					logger.Info("Row with different data detected! [%s] (%d)", rd.NodeMacAddress, rd.SequenceNumber)
					logger.Infof("  [co2] %d vs %d", tCo2, sd.co2)
					logger.Infof("  [humidity] %d vs %d", tHumidity, sd.humidity)
					logger.Infof("  [light] %d vs %d", tLight, sd.light)
					logger.Infof("  [pressure] %d vs %d", tPressure, sd.pressure)
					logger.Infof("  [noise] %d vs %d", tSoundPressure, sd.soundPressure)
					logger.Infof("  [temperature] %d vs %d", tTemperature, sd.temperature)
					logger.Infof("  [voc] %d vs %d", tTvoc, sd.tvoc)
					logger.Infof("  [uv] %d vs %d", tUv, sd.uv)
				}
			}

			_, err = stmtInsert.Exec(
				rd.RadioBusID, rd.Channel, rd.NodeMacAddress, rd.PacketType, rd.SequenceNumber,
				t.UnixNano(), rd.TimestampTz, sd.vBat, sd.vcc, sd.temperature,
				sd.humidity, sd.pressure, sd.co2, sd.tvoc, sd.light,
				sd.uv, sd.soundPressure, sd.portInput, sd.mag.string(), sd.acc.string(),
				sd.gyro.string(),
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
	// TODO: Complete this
	return true
}

func fixDatabaseEntries() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbData["DB_USER"], dbData["DB_PASS"], dbData["DB_NAME"]))

	if err != nil {
		logger.Errorf("Error connecting to DB: %s", err)
		// TODO: Figure out a good way to exit here
		return err
	}
	defer db.Close()

	rows, err := db.Query(
		"SELECT `node_mac_address`, `sequence_number`, COUNT(*) AS entry_count " +
			"FROM radio_datas " +
			"GROUP BY `node_mac_address`, `sequence_number`")

	if err != nil {
		logger.Errorf("Error querying DB: %s", err)
		return err
	}
	defer rows.Close()

	var (
		resultCountOne    int
		resultCountTwo    int
		rowsEliminatedOne int
		rowsEliminatedTwo int
	)

	for rows.Next() {
		var (
			nodeMacAddress string
			sequenceNumber int64
			entryCount     int
		)

		if err := rows.Scan(&nodeMacAddress, &sequenceNumber, &entryCount); err != nil {
			logger.Errorf("Error reading row data: %s", err)
			return err
		}

		if entryCount == 1 {
			resultCountOne++
			rowsEliminatedOne++
		}

		if entryCount > 1 {
			resultCountTwo++
			rowsEliminatedTwo += entryCount - 1
		}
	}

	logger.Infof("Found %d entries with duplicates.", resultCountTwo)
	logger.Infof("Eliminating duplicates would remove %d rows.", rowsEliminatedTwo)
	logger.Infof("Found %d entries without duplicates.", resultCountOne)
	logger.Infof("Eliminating these plus duplicates would save %d rows.", rowsEliminatedOne+rowsEliminatedTwo)

	return nil
}

func shouldDiscardData(sd *sensorDataCombined) bool {
	return sd.pressure == 0 || sd.humidity == 0 || sd.temperature == 0 || sd.co2 == 0
}
