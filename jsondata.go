package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ParseData(data string) (rd RadioData) {
	json.Unmarshal([]byte(data), &rd)
	return
}

func (v Vector3) String() string {
	return fmt.Sprintf("%d %d %d", v.X, v.Y, v.Z)
}

func (sd SensorData) GetSensorType() (st SensorType, oerr error) {
	val, err := strconv.ParseInt(sd.SensorType, 16, 64)

	if err != nil || val < 0 || val >= int64(ST_MAX_MARKER) {
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
