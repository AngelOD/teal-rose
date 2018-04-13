package main

import (
	"github.com/kardianos/service"
)

const RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"
const MdnsServiceType = "lora_server"

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
	ST_MAX_MARKER
)

var (
	dbData   = map[string]string{
		"DB_NAME": "teal_rose",
		"DB_USER": "teal_rose",
		"DB_PASS": "",
	}
	debugLog = false
	domain   = ""
	host     = ""
	port     = 3333
	saveData = false
	rdStore  = make(chan RadioData, 20)
)

var (
	env    map[string]string
	logger service.Logger
	svc    service.Service
	prg    *program
)

type Config struct {
	Name, DisplayName, Description string
}

type program struct {
	exit    chan struct{}
	service service.Service
}

//go:generate stringer -type=SensorType
type SensorType int

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
	//Timestamp   string       `json:"TimeStamp"`
	TimestampTz string       `json:"TimeStampTZ"`
	Sensors     []SensorData `json:"Sensors"`
}
