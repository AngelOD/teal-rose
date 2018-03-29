package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type SensorType int

const (
	ST_VBAT SensorType = iota
	ST_VCC
	ST_STS31_TEMP
	ST_BME680
	ST_BME280
	ST_CCS811
	ST_APDS9200
	ST_SOUNDPRESSURE
	ST_PORTINPUT
	ST_LSM9DS1TR_MAG
	ST_LSM9DS1TR_ACC
	ST_LSM9DS1TR_GYRO
	ST_COMBINED SensorType = 999
)

type SensorData struct {
	SensorType string `json:"SensorType"`
	//SensorIndex   string `json:"SensorIndex"`
	//Length        string `json:"Length"`
	VBat          int `json:"VBat"`
	Vcc           int `json:"VCC"`
	Temp1         int `json:"Temp_STS31"`
	Temp2         int `json:"Temp_BME680"`
	Temp3         int `json:"Temp_BME280"`
	Humidity1     int `json:"Humidity_BME680"`
	Humidity2     int `json:"Humidity_BME280"`
	Pressure1     int `json:"Pressure_BME680"`
	Pressure2     int `json:"Pressure_BME280"`
	Co2           int `json:"CO2"`
	Tvoc          int `json:"TVOC"`
	Light         int `json:"Light"`
	Uv            int `json:"UV"`
	SoundPressure int `json:"Soundpressure"`
	PortInput     int `json:"Port_Input"`
	MagX          int `json:"Mag_X"`
	MagY          int `json:"Mag_Y"`
	MagZ          int `json:"Mag_Z"`
	AccX          int `json:"Acc_X"`
	AccY          int `json:"Acc_Y"`
	AccZ          int `json:"Acc_Z"`
	GyroX         int `json:"Gyro_X"`
	GyroY         int `json:"Gyro_Y"`
	GyroZ         int `json:"Gyro_Z"`
}

type Vector3 struct {
	X int
	Y int
	Z int
}

type SensorDataCombined struct {
	VBat          int
	Vcc           int
	Temperature   int
	Humidity      int
	Pressure      int
	Co2           int
	Tvoc          int
	Light         int
	Uv            int
	SoundPressure int
	PortInput     int
	Mag           Vector3
	Acc           Vector3
	Gyro          Vector3
}

type RadioData struct {
	RadioBusId int `json:"radiobusid"`
	Channel    int `json:"channel"`
	//SpreadingFactor int          `json:"spreadingfactor"`
	//Rssi            int          `json:"RSSI"`
	//Snr             int          `json:"SNR"`
	NodeMacAddress string `json:"node_mac_address"`
	PacketType     int    `json:"packet_type"`
	SequenceNumber int    `json:"sequencenumber"`
	//PayloadLength   int          `json:"payloadlength"`
	//Payload         string       `json:"payload"`
	//CombinedRssiSnr float64      `json:"combined_rssi_snr"`
	Timestamp   string       `json:"TimeStamp"`
	TimestampTz string       `json:"TimeStampTZ"`
	Sensors     []SensorData `json:"Sensors"`
}

func TestJson() {
	testData := `{ "payloadlength": 92, "combined_rssi_snr": -33.75, "sequencenumber": 50, "TimeStamp": "2018-03-19 13:27:22.087292", "node_mac_address": "0000000A", "TimeStampTZ": "2018-03-19T13:27:22.087451+01:00", "radiobusid": 2, "SNR": 41, "RSSI": 113, "spreadingfactor": 8, "Sensors": [ { "SensorType": "00", "Length": "02", "VBat": 407, "SensorIndex": "00" }, { "SensorType": "01", "Length": "02", "VCC": 287, "SensorIndex": "00" }, { "SensorType": "02", "Temp_STS31": 2268, "Length": "02", "SensorIndex": "00" }, { "Humidity_BME280": 18930, "SensorType": "04", "Length": "0C", "Pressure_BME280": 101304, "Temp_BME280": 2355, "SensorIndex": "00" }, { "SensorType": "05", "Length": "0C", "CO2": 414, "TVOC": 2, "SensorIndex": "00" }, { "SensorType": "06", "Light": 8, "Length": "08", "UV": 0, "SensorIndex": "00" }, { "SensorType": "07", "Soundpressure": 66, "Length": "02", "SensorIndex": "00" }, { "SensorType": "08", "Length": "01", "Port_Input": 193, "SensorIndex": "00" }, { "Mag_X": 62376, "SensorType": "09", "Length": "06", "Mag_Y": 2842, "SensorIndex": "00", "Mag_Z": 63022 }, { "SensorIndex": "00", "SensorType": "0A", "Length": "06", "Acc_Z": 16040, "Acc_Y": 65469, "Acc_X": 65218 }, { "SensorType": "0B", "Length": "06", "Gyro_Z": 107, "Gyro_X": 34, "Gyro_Y": 130, "SensorIndex": "00" } ], "packet_type": 1, "payload": "0000020197010002011F02000208DC04000C00000933000049F200018BB805000C019E0002000000000000000006000800000008000000000700020042080001C1090006F3A80B1AF62E0A0006FEC2FFBD3EA80B000600220082006B", "channel": 0}`
	data := ParseData(testData)

	fmt.Printf("Combined: %+v\n", data.GetSensorData())
}

func ParseData(data string) (rd RadioData) {
	json.Unmarshal([]byte(data), &rd)
	return
}

func (sd SensorData) GetSensorType() (st SensorType, oerr error) {
	val, err := strconv.Atoi(sd.SensorType)

	if err != nil {
		val = -1
	}

	st = SensorType(val)
	oerr = err

	return
}

func (rd RadioData) GetSensorData() (sd SensorDataCombined) {
	var sensor SensorData

	humidityCount := 0
	pressureCount := 0
	tempCount := 0
	sensorCount := len(rd.Sensors)

	for i := 0; i < sensorCount; i++ {
		sensor = rd.Sensors[i]
		sensorType, err := sensor.GetSensorType()

		if err != nil {
			continue
		}

		switch sensorType {
		case ST_APDS9200:
			sd.Light = sensor.Light
			sd.Uv = sensor.Uv

		case ST_BME280:
			sd.Humidity += sensor.Humidity2
			sd.Pressure += sensor.Pressure2
			sd.Temperature += sensor.Temp3
			humidityCount++
			pressureCount++
			tempCount++

		case ST_BME680:
			sd.Humidity += sensor.Humidity1
			sd.Pressure += sensor.Pressure1
			sd.Temperature += sensor.Temp2
			humidityCount++
			pressureCount++
			tempCount++

		case ST_CCS811:
			sd.Co2 = sensor.Co2
			sd.Tvoc = sensor.Tvoc

		case ST_LSM9DS1TR_ACC:
			sd.Acc.X = sensor.AccX
			sd.Acc.Y = sensor.AccY
			sd.Acc.Z = sensor.AccZ

		case ST_LSM9DS1TR_GYRO:
			sd.Gyro.X = sensor.GyroX
			sd.Gyro.Y = sensor.GyroY
			sd.Gyro.Z = sensor.GyroZ

		case ST_LSM9DS1TR_MAG:
			sd.Mag.X = sensor.MagX
			sd.Mag.Y = sensor.MagY
			sd.Mag.Z = sensor.MagZ

		case ST_PORTINPUT:
			sd.PortInput = sensor.PortInput

		case ST_SOUNDPRESSURE:
			sd.SoundPressure = sensor.SoundPressure

		case ST_STS31_TEMP:
			sd.Temperature = sensor.Temp1
			tempCount++

		case ST_VBAT:
			sd.VBat = sensor.VBat

		case ST_VCC:
			sd.Vcc = sensor.Vcc
		}
	}

	if humidityCount > 1 {
		sd.Humidity = int(Round(float64(sd.Humidity) / float64(humidityCount)))
	}

	if pressureCount > 1 {
		sd.Pressure = int(Round(float64(sd.Pressure) / float64(pressureCount)))
	}

	if tempCount > 1 {
		sd.Temperature = int(Round(float64(sd.Temperature) / float64(tempCount)))
	}

	return
}
