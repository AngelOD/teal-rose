package main

import "testing"

func TestGetSensorTypeToReturnCorrectType(t *testing.T) {
	tables := []struct {
		input  string
		output sensorType
	}{
		{"00", stVBAT},
		{"01", stVCC},
		{"02", stSTS31TEMP},
		{"03", stBME680},
		{"04", stBME280},
		{"05", stCCS811},
		{"06", stAPDS9200},
		{"07", stSOUNDPRESSURE},
		{"08", stPORTINPUT},
		{"09", stLSM9DS1TRMAG},
		{"0A", stLSM9DS1TRACC},
		{"0B", stLSM9DS1TRGYRO},
		{"0", stVBAT},
		{"1", stVCC},
		{"2", stSTS31TEMP},
		{"3", stBME680},
		{"4", stBME280},
		{"5", stCCS811},
		{"6", stAPDS9200},
		{"7", stSOUNDPRESSURE},
		{"8", stPORTINPUT},
		{"9", stLSM9DS1TRMAG},
		{"A", stLSM9DS1TRACC},
		{"B", stLSM9DS1TRGYRO},
	}

	var sd sensorData

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

	var sd sensorData

	for _, testCase := range testCases {
		sd.SensorType = testCase

		st, err := sd.GetSensorType()

		if err != nil {
			t.Errorf("Got error while converting: %s", err)
		} else if st != sensorType(-1) {
			t.Errorf("Returned type was incorrect. Got %d, wanted %d.", st, -1)
		}
	}
}
