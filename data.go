package main

import (
	"github.com/kardianos/service"
)

const rfc3339Micro = "2006-01-02T15:04:05.000000Z07:00"
const mdnsServiceType = "_lora_server._tcp"

const (
	stVBAT sensorType = iota
	stVCC
	stSTS31TEMP
	stBME680
	stBME280
	stCCS811
	stAPDS9200
	stSOUNDPRESSURE
	stPORTINPUT
	stLSM9DS1TRMAG
	stLSM9DS1TRACC
	stLSM9DS1TRGYRO
	stMAXMARKER
)

var (
	dbData = map[string]string{
		"DB_TYPE": "influxdb",
		"DB_HOST": "localhost",
		"DB_PORT": "8086",
		"DB_NAME": "teal_rose",
		"DB_USER": "teal_rose",
		"DB_PASS": "",
	}
	debugLog  = false
	domain    = ""
	host      = ""
	port      = 3333
	rdStore   = make(chan radioData, 20)
	saveData  = false
	webDomain = ""
	webHost   = ""
	webIP     = ""
	webPort   = 80
)

var (
	env         map[string]string
	logger      service.Logger
	svc         service.Service
	prg         *program
	versionInfo string // Set during build phase
)

type config struct {
	name, displayName, description string
}

type program struct {
	exit    chan struct{}
	service service.Service
}

//go:generate stringer -type=sensorType
type sensorType int

type sensorData struct {
	SensorType string `json:"SensorType"`
	//sensorIndex   string `json:"SensorIndex"`
	//length        string `json:"Length"`
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

type vector3 struct {
	x int
	y int
	z int
}

type sensorDataCombined struct {
	vBat          int
	vcc           int
	temperature   int
	humidity      int
	pressure      int
	co2           int
	tvoc          int
	light         int
	uv            int
	soundPressure int
	portInput     int
	mag           vector3
	acc           vector3
	gyro          vector3
}

type radioData struct {
	RadioBusID int `json:"radiobusid"`
	Channel    int `json:"channel"`
	NodeMacAddress string `json:"node_mac_address"`
	PacketType     int    `json:"packet_type"`
	SequenceNumber int    `json:"sequencenumber"`
	TimestampTz string       `json:"TimeStampTZ"`
	Sensors     []sensorData `json:"Sensors"`
	//SpreadingFactor int          `json:"spreadingfactor"`
	//Rssi            int          `json:"RSSI"`
	//Snr             int          `json:"SNR"`
	//PayloadLength   int          `json:"payloadlength"`
	//Payload         string       `json:"payload"`
	//CombinedRssiSnr float64      `json:"combined_rssi_snr"`
	//Timestamp   string       `json:"TimeStamp"`
}
