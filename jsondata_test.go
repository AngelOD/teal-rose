package main

import "testing"

func TestGetSensorTypeToReturnCorrectType(t *testing.T) {
	tables := []struct {
		input  string
		output SensorType
	}{
		{"00", ST_VBAT},
		{"01", ST_VCC},
		{"02", ST_STS31_TEMP},
		{"03", ST_BME680},
		{"04", ST_BME280},
		{"05", ST_CCS811},
		{"06", ST_APDS9200},
		{"07", ST_SOUNDPRESSURE},
		{"08", ST_PORTINPUT},
		{"09", ST_LSM9DS1TR_MAG},
		{"0A", ST_LSM9DS1TR_ACC},
		{"0B", ST_LSM9DS1TR_GYRO},
		{"0", ST_VBAT},
		{"1", ST_VCC},
		{"2", ST_STS31_TEMP},
		{"3", ST_BME680},
		{"4", ST_BME280},
		{"5", ST_CCS811},
		{"6", ST_APDS9200},
		{"7", ST_SOUNDPRESSURE},
		{"8", ST_PORTINPUT},
		{"9", ST_LSM9DS1TR_MAG},
		{"A", ST_LSM9DS1TR_ACC},
		{"B", ST_LSM9DS1TR_GYRO},
	}

	var sd SensorData

	for _, table := range tables {
		sd.SensorType = table.input

		st, err := sd.GetSensorType()

		if err != nil {
			t.Errorf("Got error while converting: %s", err)
		} else if st != table.output {
			t.Errorf("Returned type was incorrect. Got %s, wanted %s.", st.String(), table.output.String())
		}
	}
}

func TestGetSensorTypeToReturnMinusOneIfNotFound(t *testing.T) {
	testCases := []string{
		"-01", "-1", "0C", "C", "0D", "D", "FF", "0100", "100",
	}

	var sd SensorData

	for _, testCase := range testCases {
		sd.SensorType = testCase

		st, err := sd.GetSensorType()

		if err != nil {
			t.Errorf("Got error while converting: %s", err)
		} else if st != SensorType(-1) {
			t.Errorf("Returned type was incorrect. Got %d, wanted %d.", st, -1)
		}
	}
}
