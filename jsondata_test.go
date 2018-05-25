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

		st, err := sd.getSensorType()

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

		st, err := sd.getSensorType()

		if err != nil {
			t.Errorf("Got error while converting: %s", err)
		} else if st != sensorType(-1) {
			t.Errorf("Returned type was incorrect. Got %d, wanted %d.", st, -1)
		}
	}
}

func TestParseDataToReturnCorrectlyParsedRadioData(t *testing.T) {
	tables := []struct {
		input       string
		output      radioData
		sensorCount int
	}{
		{
			input: `{"payloadlength": 53, "combined_rssi_snr": -110.5, "sequencenumber": 1413, "TimeStamp": "2018-04-17 15:05:32.330502", "node_mac_address": "000000F2", "TimeStampTZ": "2018-04-17T15:05:32.330657+02:00", "radiobusid": 1, "SNR": 6, "RSSI": 45, "spreadingfactor": 8, "Sensors": [{"SensorType": "00", "Length": "02", "VBat": 445, "SensorIndex": "00"}, {"SensorType": "01", "Length": "02", "VCC": 287, "SensorIndex": "00"}, {"SensorType": "02", "Temp_STS31": 2400, "Length": "02", "SensorIndex": "00"}, {"Humidity_BME280": 26064, "SensorType": "04", "Length": "0C", "Pressure_BME280": 102171, "Temp_BME280": 2375, "SensorIndex": "00"}, {"SensorType": "05", "Length": "04", "CO2": 473, "TVOC": 11, "SensorIndex": "00"}, {"SensorType": "06", "Light": 16, "Length": "08", "UV": 0, "SensorIndex": "00"}, {"SensorType": "07", "Soundpressure": 6, "Length": "02", "SensorIndex": "00"}], "packet_type": 1, "payload": "00000201BD010002011F020002096004000C00000947000065D000018F1B05000401D9000B06000800000010000000000700020006", "channel": 0}`,
			output: radioData{
				RadioBusID:     1,
				Channel:        0,
				NodeMacAddress: "000000F2",
				PacketType:     1,
				SequenceNumber: 1413,
				TimestampTz:    "2018-04-17T15:05:32.330657+02:00",
			},
			sensorCount: 7,
		},
		{
			input: `{"payloadlength": 3, "combined_rssi_snr": -36.75, "sequencenumber": 193, "TimeStamp": "2018-04-17 15:05:40.085026", "node_mac_address": "000000F0", "TimeStampTZ": "2018-04-17T15:05:40.085181+02:00", "radiobusid": 1, "SNR": 37, "RSSI": 111, "spreadingfactor": 8, "packet_type": 2, "payload": "01D6C3", "channel": 0}`,
			output: radioData{
				RadioBusID:     1,
				Channel:        0,
				NodeMacAddress: "000000F0",
				PacketType:     2,
				SequenceNumber: 193,
				TimestampTz:    "2018-04-17T15:05:40.085181+02:00",
			},
			sensorCount: 0,
		},
		{
			input: `{"payloadlength": 84, "combined_rssi_snr": -121.75, "sequencenumber": 1101, "TimeStamp": "2018-04-16 07:37:04.471119", "node_mac_address": "00000069", "TimeStampTZ": "2018-04-16T07:37:04.471218+02:00", "radiobusid": 2, "SNR": 225, "RSSI": 43, "spreadingfactor": 8, "Sensors": [{"SensorType": "00", "Length": "02", "VBat": 423, "SensorIndex": "00"}, {"SensorType": "01", "Length": "02", "VCC": 286, "SensorIndex": "00"}, {"SensorType": "02", "Temp_STS31": 2315, "Length": "02", "SensorIndex": "00"}, {"Humidity_BME280": 29409, "SensorType": "04", "Length": "0C", "Pressure_BME280": 101204, "Temp_BME280": 2333, "SensorIndex": "00"}, {"SensorType": "05", "Length": "04", "CO2": 413, "TVOC": 1, "SensorIndex": "00"}, {"SensorType": "06", "Light": 20, "Length": "08", "UV": 0, "SensorIndex": "00"}, {"SensorType": "07", "Soundpressure": 5, "Length": "02", "SensorIndex": "00"}, {"SensorType": "08", "Length": "01", "Port_Input": 195, "SensorIndex": "00"}, {"Mag_X": 1185, "SensorType": "09", "Length": "06", "Mag_Y": 3894, "SensorIndex": "00", "Mag_Z": 1570}, {"SensorIndex": "00", "SensorType": "0A", "Length": "06", "Acc_Z": 49631, "Acc_Y": 5034, "Acc_X": 65352}, {"SensorType": "0B", "Length": "06", "Gyro_Z": 131, "Gyro_X": 65069, "Gyro_Y": 117, "SensorIndex": "00"}], "packet_type": 1, "payload": "00000201A7010002011E020002090B04000C0000091D000072E100018B54050004019D000106000800000014000000000700020005080001C309000604A10F3606220A0006FF4813AAC1DF0B0006FE2D00750083", "channel": 0}`,
			output: radioData{
				RadioBusID:     2,
				Channel:        0,
				NodeMacAddress: "00000069",
				PacketType:     1,
				SequenceNumber: 1101,
				TimestampTz:    "2018-04-16T07:37:04.471218+02:00",
			},
			sensorCount: 11,
		},
		{
			input: `{"payloadlength": 84, "combined_rssi_snr": -117.75, "sequencenumber": 1027, "TimeStamp": "2018-04-16 01:27:01.552106", "node_mac_address": "00000069", "TimeStampTZ": "2018-04-16T01:27:01.552204+02:00", "radiobusid": 1, "SNR": 233, "RSSI": 45, "spreadingfactor": 8, "Sensors": [{"SensorType": "00", "Length": "02", "VBat": 425, "SensorIndex": "00"}, {"SensorType": "01", "Length": "02", "VCC": 286, "SensorIndex": "00"}, {"SensorType": "02", "Temp_STS31": 2326, "Length": "02", "SensorIndex": "00"}, {"Humidity_BME280": 28245, "SensorType": "04", "Length": "0C", "Pressure_BME280": 101216, "Temp_BME280": 2332, "SensorIndex": "00"}, {"SensorType": "05", "Length": "04", "CO2": 415, "TVOC": 2, "SensorIndex": "00"}, {"SensorType": "06", "Light": 20, "Length": "08", "UV": 0, "SensorIndex": "00"}, {"SensorType": "07", "Soundpressure": 5, "Length": "02", "SensorIndex": "00"}, {"SensorType": "08", "Length": "01", "Port_Input": 195, "SensorIndex": "00"}, {"Mag_X": 1158, "SensorType": "09", "Length": "06", "Mag_Y": 3935, "SensorIndex": "00", "Mag_Z": 1454}, {"SensorIndex": "00", "SensorType": "0A", "Length": "06", "Acc_Z": 49666, "Acc_Y": 5027, "Acc_X": 65355}, {"SensorType": "0B", "Length": "06", "Gyro_Z": 122, "Gyro_X": 65080, "Gyro_Y": 96, "SensorIndex": "00"}], "packet_type": 1, "payload": "00000201A9010002011E020002091604000C0000091C00006E5500018B60050004019F000206000800000014000000000700020005080001C309000604860F5F05AE0A0006FF4B13A3C2020B0006FE380060007A", "channel": 0}`,
			output: radioData{
				RadioBusID:     1,
				Channel:        0,
				NodeMacAddress: "00000069",
				PacketType:     1,
				SequenceNumber: 1027,
				TimestampTz:    "2018-04-16T01:27:01.552204+02:00",
			},
			sensorCount: 11,
		},
	}

	for _, testCase := range tables {
		rd := parseData(testCase.input)
		output := testCase.output

		if rd.RadioBusID != output.RadioBusID {
			t.Errorf("Incorrect RadioBusID. Got %d, wanted %d.", rd.RadioBusID, output.RadioBusID)
		}

		if rd.Channel != output.Channel {
			t.Errorf("Incorrect Channel. Got %d, wanted %d.", rd.Channel, output.Channel)
		}

		if rd.NodeMacAddress != output.NodeMacAddress {
			t.Errorf("Incorrect NodeMacAddress. Got %s, wanted %s.", rd.NodeMacAddress, output.NodeMacAddress)
		}

		if rd.PacketType != output.PacketType {
			t.Errorf("Incorrect PacketType. Got %d, wanted %d.", rd.PacketType, output.PacketType)
		}

		if rd.SequenceNumber != output.SequenceNumber {
			t.Errorf("Incorrect SequenceNumber. Got %d, wanted %d.", rd.SequenceNumber, output.SequenceNumber)
		}

		if rd.TimestampTz != output.TimestampTz {
			t.Errorf("Incorrect TimestampTz. Got %s, wanted %s.", rd.TimestampTz, output.TimestampTz)
		}

		if len(rd.Sensors) != testCase.sensorCount {
			t.Errorf("Incorrect sensor count. Got %d, wanted %d.", len(rd.Sensors), testCase.sensorCount)
		}
	}
}

func TestParseDataToReturnCorrectlyParsedSensorData(t *testing.T) {
	tables := []struct {
		input  string
		output []sensorData
	}{
		{
			input: `{"payloadlength": 53, "combined_rssi_snr": -110.5, "sequencenumber": 1413, "TimeStamp": "2018-04-17 15:05:32.330502", "node_mac_address": "000000F2", "TimeStampTZ": "2018-04-17T15:05:32.330657+02:00", "radiobusid": 1, "SNR": 6, "RSSI": 45, "spreadingfactor": 8, "Sensors": [{"SensorType": "00", "Length": "02", "VBat": 445, "SensorIndex": "00"}, {"SensorType": "01", "Length": "02", "VCC": 287, "SensorIndex": "00"}, {"SensorType": "02", "Temp_STS31": 2400, "Length": "02", "SensorIndex": "00"}, {"Humidity_BME280": 26064, "SensorType": "04", "Length": "0C", "Pressure_BME280": 102171, "Temp_BME280": 2375, "SensorIndex": "00"}, {"SensorType": "05", "Length": "04", "CO2": 473, "TVOC": 11, "SensorIndex": "00"}, {"SensorType": "06", "Light": 16, "Length": "08", "UV": 0, "SensorIndex": "00"}, {"SensorType": "07", "Soundpressure": 6, "Length": "02", "SensorIndex": "00"}], "packet_type": 1, "payload": "00000201BD010002011F020002096004000C00000947000065D000018F1B05000401D9000B06000800000010000000000700020006", "channel": 0}`,
			output: []sensorData{
				{
					SensorType: "00",
					VBat:       445,
				},
				{
					SensorType: "01",
					Vcc:        287,
				},
				{
					SensorType: "02",
					Temp1:      2400,
				},
				{
					SensorType: "04",
					Humidity2:  26064,
					Pressure2:  102171,
					Temp3:      2375,
				},
				{
					SensorType: "05",
					Co2:        473,
					Tvoc:       11,
				},
				{
					SensorType: "06",
					Light:      16,
					Uv:         0,
				},
				{
					SensorType:    "07",
					SoundPressure: 6,
				},
			},
		},
		{
			input: `{"payloadlength": 84, "combined_rssi_snr": -121.75, "sequencenumber": 1101, "TimeStamp": "2018-04-16 07:37:04.471119", "node_mac_address": "00000069", "TimeStampTZ": "2018-04-16T07:37:04.471218+02:00", "radiobusid": 2, "SNR": 225, "RSSI": 43, "spreadingfactor": 8, "Sensors": [{"SensorType": "00", "Length": "02", "VBat": 423, "SensorIndex": "00"}, {"SensorType": "01", "Length": "02", "VCC": 286, "SensorIndex": "00"}, {"SensorType": "02", "Temp_STS31": 2315, "Length": "02", "SensorIndex": "00"}, {"Humidity_BME280": 29409, "SensorType": "04", "Length": "0C", "Pressure_BME280": 101204, "Temp_BME280": 2333, "SensorIndex": "00"}, {"SensorType": "05", "Length": "04", "CO2": 413, "TVOC": 1, "SensorIndex": "00"}, {"SensorType": "06", "Light": 20, "Length": "08", "UV": 0, "SensorIndex": "00"}, {"SensorType": "07", "Soundpressure": 5, "Length": "02", "SensorIndex": "00"}, {"SensorType": "08", "Length": "01", "Port_Input": 195, "SensorIndex": "00"}, {"Mag_X": 1185, "SensorType": "09", "Length": "06", "Mag_Y": 3894, "SensorIndex": "00", "Mag_Z": 1570}, {"SensorIndex": "00", "SensorType": "0A", "Length": "06", "Acc_Z": 49631, "Acc_Y": 5034, "Acc_X": 65352}, {"SensorType": "0B", "Length": "06", "Gyro_Z": 131, "Gyro_X": 65069, "Gyro_Y": 117, "SensorIndex": "00"}], "packet_type": 1, "payload": "00000201A7010002011E020002090B04000C0000091D000072E100018B54050004019D000106000800000014000000000700020005080001C309000604A10F3606220A0006FF4813AAC1DF0B0006FE2D00750083", "channel": 0}`,
			output: []sensorData{
				{
					SensorType: "00",
					VBat:       423,
				},
				{
					SensorType: "01",
					Vcc:        286,
				},
				{
					SensorType: "02",
					Temp1:      2315,
				},
				{
					SensorType: "04",
					Humidity2:  29409,
					Pressure2:  101204,
					Temp3:      2333,
				},
				{
					SensorType: "05",
					Co2:        413,
					Tvoc:       1,
				},
				{
					SensorType: "06",
					Light:      20,
					Uv:         0,
				},
				{
					SensorType:    "07",
					SoundPressure: 5,
				},
				{
					SensorType: "08",
					PortInput:  195,
				},
				{
					SensorType: "09",
					MagX:       1185,
					MagY:       3894,
					MagZ:       1570,
				},
				{
					SensorType: "0A",
					AccX:       65352,
					AccY:       5034,
					AccZ:       49631,
				},
				{
					SensorType: "0B",
					GyroX:      65069,
					GyroY:      117,
					GyroZ:      131,
				},
			},
		},
		{
			input: `{"payloadlength": 84, "combined_rssi_snr": -117.75, "sequencenumber": 1027, "TimeStamp": "2018-04-16 01:27:01.552106", "node_mac_address": "00000069", "TimeStampTZ": "2018-04-16T01:27:01.552204+02:00", "radiobusid": 1, "SNR": 233, "RSSI": 45, "spreadingfactor": 8, "Sensors": [{"SensorType": "00", "Length": "02", "VBat": 425, "SensorIndex": "00"}, {"SensorType": "01", "Length": "02", "VCC": 286, "SensorIndex": "00"}, {"SensorType": "02", "Temp_STS31": 2326, "Length": "02", "SensorIndex": "00"}, {"Humidity_BME280": 28245, "SensorType": "04", "Length": "0C", "Pressure_BME280": 101216, "Temp_BME280": 2332, "SensorIndex": "00"}, {"SensorType": "05", "Length": "04", "CO2": 415, "TVOC": 2, "SensorIndex": "00"}, {"SensorType": "06", "Light": 20, "Length": "08", "UV": 0, "SensorIndex": "00"}, {"SensorType": "07", "Soundpressure": 5, "Length": "02", "SensorIndex": "00"}, {"SensorType": "08", "Length": "01", "Port_Input": 195, "SensorIndex": "00"}, {"Mag_X": 1158, "SensorType": "09", "Length": "06", "Mag_Y": 3935, "SensorIndex": "00", "Mag_Z": 1454}, {"SensorIndex": "00", "SensorType": "0A", "Length": "06", "Acc_Z": 49666, "Acc_Y": 5027, "Acc_X": 65355}, {"SensorType": "0B", "Length": "06", "Gyro_Z": 122, "Gyro_X": 65080, "Gyro_Y": 96, "SensorIndex": "00"}], "packet_type": 1, "payload": "00000201A9010002011E020002091604000C0000091C00006E5500018B60050004019F000206000800000014000000000700020005080001C309000604860F5F05AE0A0006FF4B13A3C2020B0006FE380060007A", "channel": 0}`,
			output: []sensorData{
				{
					SensorType: "00",
					VBat:       425,
				},
				{
					SensorType: "01",
					Vcc:        286,
				},
				{
					SensorType: "02",
					Temp1:      2326,
				},
				{
					SensorType: "04",
					Humidity2:  28245,
					Pressure2:  101216,
					Temp3:      2332,
				},
				{
					SensorType: "05",
					Co2:        415,
					Tvoc:       2,
				},
				{
					SensorType: "06",
					Light:      20,
					Uv:         0,
				},
				{
					SensorType:    "07",
					SoundPressure: 5,
				},
				{
					SensorType: "08",
					PortInput:  195,
				},
				{
					SensorType: "09",
					MagX:       1158,
					MagY:       3935,
					MagZ:       1454,
				},
				{
					SensorType: "0A",
					AccX:       65355,
					AccY:       5027,
					AccZ:       49666,
				},
				{
					SensorType: "0B",
					GyroX:      65080,
					GyroY:      96,
					GyroZ:      122,
				},
			},
		},
	}

	for _, testCase := range tables {
		rd := parseData(testCase.input)

		for _, sensor := range testCase.output {
			found := false

			for _, fSensor := range rd.Sensors {
				if AssertThatSensorDataStructsAreEqual(&sensor, &fSensor) {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Incorrect sensor data.\nWanted:\n%#v\nGot:\n%#v", testCase.output, rd.Sensors)
			}
		}
	}
}

func BenchmarkParseData(b *testing.B) {
	testData := `{"payloadlength": 53, "combined_rssi_snr": -110.5, "sequencenumber": 1413, "TimeStamp": "2018-04-17 15:05:32.330502", "node_mac_address": "000000F2", "TimeStampTZ": "2018-04-17T15:05:32.330657+02:00", "radiobusid": 1, "SNR": 6, "RSSI": 45, "spreadingfactor": 8, "Sensors": [{"SensorType": "00", "Length": "02", "VBat": 445, "SensorIndex": "00"}, {"SensorType": "01", "Length": "02", "VCC": 287, "SensorIndex": "00"}, {"SensorType": "02", "Temp_STS31": 2400, "Length": "02", "SensorIndex": "00"}, {"Humidity_BME280": 26064, "SensorType": "04", "Length": "0C", "Pressure_BME280": 102171, "Temp_BME280": 2375, "SensorIndex": "00"}, {"SensorType": "05", "Length": "04", "CO2": 473, "TVOC": 11, "SensorIndex": "00"}, {"SensorType": "06", "Light": 16, "Length": "08", "UV": 0, "SensorIndex": "00"}, {"SensorType": "07", "Soundpressure": 6, "Length": "02", "SensorIndex": "00"}], "packet_type": 1, "payload": "00000201BD010002011F020002096004000C00000947000065D000018F1B05000401D9000B06000800000010000000000700020006", "channel": 0}`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parseData(testData)
	}
}

func AssertThatSensorDataStructsAreEqual(a *sensorData, b *sensorData) bool {
	return a.SensorType == b.SensorType &&
		a.VBat == b.VBat &&
		a.Vcc == b.Vcc &&
		a.Temp1 == b.Temp1 &&
		a.Temp2 == b.Temp2 &&
		a.Temp3 == b.Temp3 &&
		a.Humidity1 == b.Humidity1 &&
		a.Humidity2 == b.Humidity2 &&
		a.Pressure1 == b.Pressure1 &&
		a.Pressure2 == b.Pressure2 &&
		a.Co2 == b.Co2 &&
		a.Tvoc == b.Tvoc &&
		a.Light == b.Light &&
		a.Uv == b.Uv &&
		a.SoundPressure == b.SoundPressure &&
		a.PortInput == b.PortInput &&
		a.MagX == b.MagX &&
		a.MagY == b.MagY &&
		a.MagZ == b.MagZ &&
		a.AccX == b.AccX &&
		a.AccY == b.AccY &&
		a.AccZ == b.AccZ &&
		a.GyroX == b.GyroX &&
		a.GyroY == b.GyroY &&
		a.GyroZ == b.GyroZ
}
