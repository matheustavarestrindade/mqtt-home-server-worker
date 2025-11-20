package water_meter_worker

import (
	"fmt"
	wm_payload_v1 "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/water_meter/payloads/water_meter_payload_v1"
)

type SeverityLevel int

const (
	NORMAL SeverityLevel = iota
	WARNING
	CRITICAL
)

type WaterLevelMeterSensorData struct {
	AverageWaterLevelCm float32 `json:"average_water_level_cm"`
}

type WaterLevelMeterSensorDataResponse struct {
	PayloadVersion int `json:"v"`
	WaterLevelMeterSensorData
}

func CalculateSeverityLevel(value int) SeverityLevel {
	switch {
	case value >= 80:
		return CRITICAL
	case value >= 50:
		return WARNING
	default:
		return NORMAL
	}
}

func ConvertCompressedPayloadToSensorDataResponse(payloadVersion int, payload string) *WaterLevelMeterSensorDataResponse {
	switch payloadVersion {
	case WaterMeterMessageV1:
		data, err := wm_payload_v1.DecompressDataFromDatabase(payload)
		if err != nil {
			fmt.Printf("Failed to decompress Hydroponic Manager v1 data: %v\n", err)
			break
		}

		return &WaterLevelMeterSensorDataResponse{
			PayloadVersion: WaterMeterMessageV1,
			WaterLevelMeterSensorData: WaterLevelMeterSensorData{
				AverageWaterLevelCm: data.Sensors.AverageWaterLevelCm,
			},
		}
	default:
		fmt.Printf("Unsupported Hydroponic Manager message version: %d\n", payloadVersion)
		return nil
	}

	return &WaterLevelMeterSensorDataResponse{
		PayloadVersion: WaterMeterMessageV1,
	}
}
