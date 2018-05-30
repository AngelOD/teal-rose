package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

const testCaseCount = 10000

func setupMysqlStoreDataRunnerTest(t *testing.T) (db *sql.DB) {
	initService()
	prg.exit = make(chan struct{})

	if !loadDotEnv(".test.env") {
		dbData["DB_USER"] = "teal_rose"
		dbData["DB_PASS"] = ""
		dbData["DB_HOST"] = "localhost"
		dbData["DB_PORT"] = "3306"
		dbData["DB_NAME"] = "teal_rose_test"
	} else {
		if retVal := parseDbCli(nil, nil); retVal != 0 {
			t.Errorf("Parsing returned error code: %d", retVal)
			return
		}
	}

	fmt.Printf("%+v\n", dbData)

	conf := mysql.Config{
		User:   dbData["DB_USER"],
		Passwd: dbData["DB_PASS"],
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", dbData["DB_HOST"], dbData["DB_PORT"]),
		DBName: dbData["DB_NAME"],
	}

	db, err := sql.Open("mysql", conf.FormatDSN())

	if err != nil {
		t.Errorf("Error connecting to DB: %s", err)
		return
	}

	db.SetMaxOpenConns(1000)

	if _, err := db.Exec("TRUNCATE radio_datas"); err != nil {
		t.Errorf("Error executing TRUNCATE: %s", err)
		return
	}

	return
}

func TestMysqlStoreDataRunnerSavesData(t *testing.T) {
	db := setupMysqlStoreDataRunnerTest(t)
	defer db.Close()

	data := make([]radioData, 0, testCaseCount)

	go mysqlStoreDataRunner()

	for i := 0; i < testCaseCount; i++ {
		entry := makeRandomRadioDataForTesting()
		data = append(data, entry)
		rdStore <- entry
	}

	//sleepyTime, _ := time.ParseDuration(fmt.Sprintf("%ds", int(math.Ceil(testCaseCount/250))))
	//time.Sleep(sleepyTime)
	for len(rdStore) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	close(prg.exit)
	time.Sleep(500 * time.Millisecond)

	// Verify number first
	cRows, err := db.Query("SELECT COUNT(*) as count FROM radio_datas")
	if err != nil {
		t.Errorf("Unable to query count: %s", err)
	} else {
		var count int
		cRows.Next()
		cRows.Scan(&count)

		if count != testCaseCount {
			t.Errorf("Found %d rows, expected %d!", count, testCaseCount)
		}
	}
	cRows.Close()

	for _, entry := range data {
		var id int
		sd := entry.getSensorData()
		query :=
			"SELECT `id` FROM radio_datas" +
				" WHERE `radio_bus_id`=? AND `channel`=? AND `node_mac_address`=? AND `packet_type`=?" +
				" AND `sequence_number`=? AND `timestamp_tz`=? AND `v_bat`=?" +
				" AND `vcc`=? AND `temperature`=? AND `humidity`=? AND `pressure`=? AND `co2`=?" +
				" AND `tvoc`=? AND `light`=? AND `uv`=? AND `sound_pressure`=? AND `port_input`=?" +
				" AND `mag`=? AND `acc`=? AND `gyro`=?"
		err := db.QueryRow(
			query,
			entry.RadioBusID, entry.Channel, entry.NodeMacAddress, entry.PacketType,
			entry.SequenceNumber, entry.TimestampTz, sd.vBat,
			sd.vcc, sd.temperature, sd.humidity, sd.pressure, sd.co2,
			sd.tvoc, sd.light, sd.uv, sd.soundPressure, sd.portInput,
			sd.mag.string(), sd.acc.string(), sd.gyro.string(),
		).Scan(&id)

		if err != nil {
			t.Errorf("Row not found! (%s)\n%#v", err, entry)
			continue
		}
	}
}

func makeRandomRadioDataForTesting() radioData {
	rd := radioData{
		Channel:        rand.Intn(11),
		NodeMacAddress: strconv.FormatInt(int64(rand.Intn(2000000000)), 16),
		PacketType:     rand.Intn(2) + 1,
		RadioBusID:     rand.Intn(2),
		SequenceNumber: rand.Intn(10000),
		TimestampTz: fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%06dZ",
			rand.Intn(2)+2017, rand.Intn(12)+1, rand.Intn(28)+1,
			rand.Intn(24), rand.Intn(60), rand.Intn(60), rand.Intn(1000000)),
	}

	rd.Sensors = []sensorData{
		{
			SensorType: strconv.FormatInt(int64(stVBAT), 16),
			VBat:       rand.Intn(100) + 250,
		},
		{
			SensorType: strconv.FormatInt(int64(stVCC), 16),
			Vcc:        rand.Intn(200) + 350,
		},
		{
			SensorType: strconv.FormatInt(int64(stSTS31TEMP), 16),
			Temp1:      rand.Intn(13000) + 15000,
		},
		{
			SensorType: strconv.FormatInt(int64(stBME280), 16),
			Humidity2:  rand.Intn(4100) + 4000,
			Pressure2:  rand.Intn(100) + 101200,
			Temp3:      rand.Intn(13000) + 15000,
		},
		{
			SensorType: strconv.FormatInt(int64(stAPDS9200), 16),
			Light:      rand.Intn(1000),
			Uv:         rand.Intn(12),
		},
		{
			SensorType: strconv.FormatInt(int64(stCCS811), 16),
			Co2:        rand.Intn(7600) + 400,
			Tvoc:       rand.Intn(1100),
		},
		{
			SensorType: strconv.FormatInt(int64(stPORTINPUT), 16),
			PortInput:  rand.Intn(200) + 1,
		},
		{
			SensorType:    strconv.FormatInt(int64(stSOUNDPRESSURE), 16),
			SoundPressure: rand.Intn(100) + 5,
		},
		{
			SensorType: strconv.FormatInt(int64(stLSM9DS1TRACC), 16),
			AccX:       rand.Intn(10000),
			AccY:       rand.Intn(10000),
			AccZ:       rand.Intn(10000),
		},
		{
			SensorType: strconv.FormatInt(int64(stLSM9DS1TRGYRO), 16),
			GyroX:      rand.Intn(10000),
			GyroY:      rand.Intn(10000),
			GyroZ:      rand.Intn(10000),
		},
		{
			SensorType: strconv.FormatInt(int64(stLSM9DS1TRMAG), 16),
			MagX:       rand.Intn(10000),
			MagY:       rand.Intn(10000),
			MagZ:       rand.Intn(10000),
		},
	}

	return rd
}
