package water_meter_payload_v1

import (
	"fmt"
	"strconv"
	"strings"
)

type Payload struct {
	Version  int    `json:"version"`
	ClientId string `json:"clientId"`
	FuseId   string `json:"espFuseId"`
	Data     Data   `json:"data"`
}

type Data struct {
	Sensors SensorData `json:"sensors"`
}

type SeverityLevel int

const MAX_COMPRESSED_PAYLOAD_LENGTH = 128

const (
	NORMAL SeverityLevel = iota
	WARNING
	CRITICAL
)

type SensorData struct {
	AverageWaterLevelCm float32 `json:"averageWaterLevelCm"`
}

func ParsePayload(parts []string) (*Payload, error) {
	var message Payload

	message.ClientId = parts[1]
	message.FuseId = parts[2]

	for i := 3; i < len(parts); i++ {
		var values = strings.Split(parts[i], ":")
		var key = values[0]

		if key == "wl" {
			if len(values) != 2 {
				return nil, fmt.Errorf("invalid water level data format")
			}
			waterLevel, err := strconv.ParseFloat(values[1], 32)
			if err != nil {
				return nil, fmt.Errorf("invalid water level value: %s", values[1])
			}
			message.Data.Sensors.AverageWaterLevelCm = float32(waterLevel)
			continue
		}
	}

	return &message, nil
}

func CompressDataToDatabase(data Data) (string, error) {
	const compressedDataFormat = "wl:%.2f"
	compressedData := fmt.Sprintf(compressedDataFormat, data.Sensors.AverageWaterLevelCm)
	if len(compressedData) > MAX_COMPRESSED_PAYLOAD_LENGTH {
		fmt.Printf("Warning: Compressed data length %d exceeds maximum of %d characters\n", len(compressedData), MAX_COMPRESSED_PAYLOAD_LENGTH)
		return "", fmt.Errorf("compressed data exceeds maximum length of %d characters", MAX_COMPRESSED_PAYLOAD_LENGTH)
	}
	return compressedData, nil
}

func DecompressDataFromDatabase(compressedData string) (Data, error) {
	var data Data

	for part := range strings.SplitSeq(compressedData, ";") {
		values := strings.Split(part, ":")
		key := values[0]
		switch key {
		case "wl":
			if len(values) != 2 {
				return data, fmt.Errorf("invalid water level data format")
			}
			waterLevel, err := strconv.ParseFloat(values[1], 32)
			if err != nil {
				return data, fmt.Errorf("invalid water level value: %s", values[1])
			}
			data.Sensors.AverageWaterLevelCm = float32(waterLevel)
		default:
			fmt.Printf("Warning: Unknown data key: %s\n", key)
			fmt.Printf("Values for unknown key: %v\n", values)
			return data, fmt.Errorf("unknown data key: %s", key)
		}
	}

	return data, nil
}
