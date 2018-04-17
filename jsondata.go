package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ParseData(data string) (rd radioData) {
	json.Unmarshal([]byte(data), &rd)
	return
}

func (v vector3) String() string {
	return fmt.Sprintf("%d %d %d", v.x, v.y, v.z)
}

func (sd sensorData) GetSensorType() (st sensorType, oerr error) {
	val, err := strconv.ParseInt(sd.SensorType, 16, 64)

	if err != nil || val < 0 || val >= int64(stMAXMARKER) {
		val = -1
	}

	st = sensorType(val)
	oerr = err

	return
}

func (rd radioData) GetSensorData() (sd sensorDataCombined) {
	var sensor sensorData

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
		case stAPDS9200:
			sd.light = sensor.Light
			sd.uv = sensor.Uv

		case stBME280:
			sd.humidity += sensor.Humidity2
			sd.pressure += sensor.Pressure2
			sd.temperature += sensor.Temp3
			humidityCount++
			pressureCount++
			tempCount++

		case stBME680:
			sd.humidity += sensor.Humidity1
			sd.pressure += sensor.Pressure1
			sd.temperature += sensor.Temp2
			humidityCount++
			pressureCount++
			tempCount++

		case stCCS811:
			sd.co2 = sensor.Co2
			sd.tvoc = sensor.Tvoc

		case stLSM9DS1TRACC:
			sd.acc.x = sensor.AccX
			sd.acc.y = sensor.AccY
			sd.acc.z = sensor.AccZ

		case stLSM9DS1TRGYRO:
			sd.gyro.x = sensor.GyroX
			sd.gyro.y = sensor.GyroY
			sd.gyro.z = sensor.GyroZ

		case stLSM9DS1TRMAG:
			sd.mag.x = sensor.MagX
			sd.mag.y = sensor.MagY
			sd.mag.z = sensor.MagZ

		case stPORTINPUT:
			sd.portInput = sensor.PortInput

		case stSOUNDPRESSURE:
			sd.soundPressure = sensor.SoundPressure

		case stSTS31TEMP:
			sd.temperature = sensor.Temp1
			tempCount++

		case stVBAT:
			sd.vBat = sensor.VBat

		case stVCC:
			sd.vcc = sensor.Vcc
		}
	}

	if humidityCount > 1 {
		sd.humidity = int(round(float64(sd.humidity) / float64(humidityCount)))
	}

	if pressureCount > 1 {
		sd.pressure = int(round(float64(sd.pressure) / float64(pressureCount)))
	}

	if tempCount > 1 {
		sd.temperature = int(round(float64(sd.temperature) / float64(tempCount)))
	}

	return
}
