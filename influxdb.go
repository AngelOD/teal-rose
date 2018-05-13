package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/influxdata/influxdb/client/v2"
)

func migrateMysqlToInflux(mysqlConnector, influxHost, influxDbName, influxUser, influxPass string) {
	logger.Info("Starting data storage thread.")

	db, err := sql.Open("mysql", mysqlConnector)

	if err != nil {
		logger.Errorf("Error connecting to DB: %s", err)
		return
	}
	defer db.Close()

	stmtSelect, err := db.Prepare(
		"SELECT `id`, `radio_bus_id`, `channel`, `node_mac_address`, `packet_type`, `sequence_number`, " +
			"`timestamp_tz`, `v_bat`, `vcc`, `co2`, `humidity`, `light`, `pressure`, `sound_pressure`, " +
			"`temperature`, `tvoc`, `uv`, `port_input`, `mag`, `acc`, `gyro` " +
			"FROM radio_datas " +
			"ORDER BY id ASC " +
			"LIMIT ? " +
			"OFFSET ?")
	defer stmtSelect.Close()

	if err != nil {
		logger.Errorf("Error preparing statement: %s", err)
		return
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxHost,
		Username: influxUser,
		Password: influxPass,
	})
	if err != nil {
		logger.Errorf("Error creating HTTP Client: %s", err)
		return
	}
	defer c.Close()

	q := client.Query{
		Command:  "DELETE FROM radio_datas",
		Database: influxDbName,
	}

	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			logger.Errorf("Unable to clear InfluxDB measurement: %s", response.Error())
			return
		}
	} else {
		logger.Errorf("Unable to clear InfluxDB measurement: %s", err)
		return
	}

	var (
		limit           = 250
		offset          int
		tAcc            string
		tChannel        int
		tCo2            int
		tGyro           string
		tHumidity       int
		tID             int
		tLight          int
		tMag            string
		tNodeMacAddress string
		tPacketType     int
		tPortInput      int
		tPressure       int
		tRadioBusID     int
		tSequenceNumber int
		tSoundPressure  int
		tTemperature    int
		tTimestampTz    string
		tTvoc           int
		tUv             int
		tVBat           int
		tVcc            int
	)

	for {
		var rowCount = 0

		logger.Infof("Retrieving %d records, starting at %d.", limit, offset)
		rows, err := stmtSelect.Query(limit, offset)

		if err != nil {
			logger.Errorf("Error querying: %s", err)
			break
		}

		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  influxDbName,
			Precision: "ns",
		})
		if err != nil {
			log.Fatal(err)
			return
		}

		for rows.Next() {
			rowCount++

			err = rows.Scan(
				&tID, &tRadioBusID, &tChannel, &tNodeMacAddress, &tPacketType, &tSequenceNumber,
				&tTimestampTz, &tVBat, &tVcc, &tCo2, &tHumidity, &tLight, &tPressure, &tSoundPressure,
				&tTemperature, &tTvoc, &tUv, &tPortInput, &tMag, &tAcc, &tGyro,
			)

			if err != nil {
				logger.Errorf("Error parsing data: %s", err)
				break
			}

			t, err := time.Parse(rfc3339Micro, tTimestampTz)
			if err != nil {
				logger.Warningf("Unable to convert time string: %s", tTimestampTz)
				logger.Warningf("Because of: %s", err)
				continue
			}

			tags := map[string]string{
				"node_mac_address": tNodeMacAddress,
				"sequence_number":  strconv.FormatInt(int64(tSequenceNumber), 10),
			}
			fields := map[string]interface{}{
				"radio_bus_id":   tRadioBusID,
				"channel":        tChannel,
				"packet_type":    tPacketType,
				"timestamp_tz":   tTimestampTz,
				"v_bat":          tVBat,
				"vcc":            tVcc,
				"temperature":    tTemperature,
				"humidity":       tHumidity,
				"pressure":       tPressure,
				"co2":            tCo2,
				"tvoc":           tTvoc,
				"light":          tLight,
				"uv":             tUv,
				"sound_pressure": tSoundPressure,
				"port_input":     tPortInput,
				"mag":            tMag,
				"acc":            tAcc,
				"gyro":           tGyro,
			}

			pt, err := client.NewPoint("radio_datas", tags, fields, t)
			if err != nil {
				log.Fatal(err)
				return
			}
			bp.AddPoint(pt)
		}

		if err := c.Write(bp); err != nil {
			log.Fatal(err)
			return
		}

		if rowCount < limit {
			logger.Infof("%d < %d -- Exiting", rowCount, limit)
			break
		}

		offset += limit
	}

	logger.Info("Running test query #1")

	qs := client.Query{
		Command:  "SELECT COUNT(temperature) FROM radio_datas",
		Database: influxDbName,
	}

	if response, err := c.Query(qs); err == nil {
		if response.Error() != nil {
			logger.Errorf("Unable to execute test query: %s", response.Error())
			return
		}

		count := response.Results[0].Series[0].Values[0][1]
		logger.Infof("Number of entries: %v", count)
	} else {
		logger.Errorf("Unable to execute test query: %s", err)
		return
	}

	logger.Info("Running test query #2")

	qs = client.Query{
		Command:  "SELECT temperature, humidity, co2, tvoc FROM radio_datas WHERE node_mac_address = '0000000A' LIMIT 20",
		Database: influxDbName,
	}

	if response, err := c.Query(qs); err == nil {
		if response.Error() != nil {
			logger.Errorf("Unable to execute test query: %s", response.Error())
			return
		}

		logger.Infof("Query result: %+v", response.Results)
	} else {
		logger.Errorf("Unable to execute test query: %s", err)
		return
	}
}

func influxStoreDataRunner() {
	var rd radioData
	var sd sensorDataCombined

	logger.Info("Starting InfluxDB data storage thread.")

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%s", dbData["DB_HOST"], dbData["DB_PORT"]),
		Username: dbData["DB_USER"],
		Password: dbData["DB_PASS"],
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer c.Close()

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

			q := fmt.Sprintf(
				"SELECT co2, humidity, light, pressure, sound_pressure, temperature, tvoc, uv "+
					"FROM radio_datas WHERE node_mac_address='%s' AND sequence_number='%d'",
				rd.NodeMacAddress, rd.SequenceNumber)
			res, err := queryInfluxDB(c, q)
			if err != nil {
				logger.Warningf("Error querying data: %s", err)
				continue
			}

			if len(res[0].Series) > 0 && len(res[0].Series[0].Values) > 0 {
				row := res[0].Series[0].Values[0]

				if row[1] == sd.co2 && row[2] == sd.humidity && row[3] == sd.light && row[4] == sd.pressure &&
					row[5] == sd.soundPressure && row[6] == sd.temperature && row[7] == sd.tvoc && row[8] == sd.uv {
					//
					logger.Info("Multiple entries detected. Doing nothing for now.")
				} else {
					logger.Info("Row with different data detected! [%s] (%d)", rd.NodeMacAddress, rd.SequenceNumber)
					logger.Infof("  [co2] %d vs %d", row[1], sd.co2)
					logger.Infof("  [humidity] %d vs %d", row[2], sd.humidity)
					logger.Infof("  [light] %d vs %d", row[3], sd.light)
					logger.Infof("  [pressure] %d vs %d", row[4], sd.pressure)
					logger.Infof("  [noise] %d vs %d", row[5], sd.soundPressure)
					logger.Infof("  [temperature] %d vs %d", row[6], sd.temperature)
					logger.Infof("  [voc] %d vs %d", row[7], sd.tvoc)
					logger.Infof("  [uv] %d vs %d", row[8], sd.uv)
				}
			}

			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  dbData["DB_NAME"],
				Precision: "ns",
			})
			if err != nil {
				logger.Warningf("Error creating BatchPoints: %s", err)
				continue
			}

			tags := map[string]string{
				"node_mac_address": "",
				"sequence_number":  "",
			}
			fields := map[string]interface{}{
				"radio_bus_id":   0,
				"channel":        0,
				"packet_type":    0,
				"timestamp_tz":   "",
				"v_bat":          0,
				"vcc":            0,
				"temperature":    0,
				"humidity":       0,
				"pressure":       0,
				"co2":            0,
				"tvoc":           0,
				"light":          0,
				"uv":             0,
				"sound_pressure": 0,
				"port_input":     0,
				"mag":            0,
				"acc":            0,
				"gyro":           0,
			}

			pt, err := client.NewPoint("radio_datas", tags, fields, t)
			if err != nil {
				logger.Warningf("Error creating Point: %s", err)
				continue
			}
			bp.AddPoint(pt)

			if err := c.Write(bp); err != nil {
				logger.Warningf("Error writing to DB: %s", err)
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

func queryInfluxDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: dbData["DB_NAME"],
	}

	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}

	return res, nil
}
