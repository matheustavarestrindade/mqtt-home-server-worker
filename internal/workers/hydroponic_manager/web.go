package hydroponic_manager_worker

import (
	"fmt"
	hm_payload_v1 "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/hydroponic_manager/payloads/hydroponic_manager_payload_v1"
)

type SeverityLevel int

const (
	NORMAL SeverityLevel = iota
	WARNING
	CRITICAL
)

type HydroponicManagerSensorData struct {
	Temperature          float32       `json:"temperature"`
	TemperaturaSeverity  SeverityLevel `json:"temperatureSeverity"`
	Moisture             float32       `json:"moisture"`
	MoistureSeverity     SeverityLevel `json:"moistureSeverity"`
	Ph                   float32       `json:"ph"`
	PhSeverity           SeverityLevel `json:"phSeverity"`
	Conductivity         int           `json:"conductivity"`
	ConductivitySeverity SeverityLevel `json:"conductivitySeverity"`
	Nitrogen             int           `json:"nitrogen"`
	NitrogenSeverity     SeverityLevel `json:"nitrogenSeverity"`
	Phosphorus           int           `json:"phosphorus"`
	PhosphorusSeverity   SeverityLevel `json:"phosphorusSeverity"`
	Potassium            int           `json:"potassium"`
	PotassiumSeverity    SeverityLevel `json:"potassiumSeverity"`
}

type HydroponicManagerRelay struct {
	IsOn                bool `json:"isOn"`
	NextToggleInSeconds int  `json:"nextToggleInSeconds"`
}

type HydroponicManagerSensorDataResponse struct {
	PayloadVersion int `json:"v"`
	HydroponicManagerSensorData
	HydroponicManagerRelay
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

func ConvertCompressedPayloadToSensorDataResponse(payloadVersion int, payload string) *HydroponicManagerSensorDataResponse {

	switch payloadVersion {
	case HydroponicManagerMessageV1:
		data, err := hm_payload_v1.DecompressDataFromDatabase(payload)
		if err != nil {
			fmt.Printf("Failed to decompress Hydroponic Manager v1 data: %v\n", err)
			break
		}

		sensor := data.Sensors
		relay := data.Relay

		return &HydroponicManagerSensorDataResponse{
			PayloadVersion: HydroponicManagerMessageV1,
			HydroponicManagerSensorData: HydroponicManagerSensorData{
				Temperature:          sensor.Temperature,
				TemperaturaSeverity:  SeverityLevel(sensor.TemperaturaSeverity),
				Moisture:             sensor.Moisture,
				MoistureSeverity:     SeverityLevel(sensor.MoistureSeverity),
				Ph:                   sensor.Ph,
				PhSeverity:           SeverityLevel(sensor.PhSeverity),
				Conductivity:         sensor.Conductivity,
				ConductivitySeverity: SeverityLevel(sensor.ConductivitySeverity),
				Nitrogen:             sensor.Nitrogen,
				NitrogenSeverity:     SeverityLevel(sensor.NitrogenSeverity),
				Phosphorus:           sensor.Phosphorus,
				PhosphorusSeverity:   SeverityLevel(sensor.PhosphorusSeverity),
				Potassium:            sensor.Potassium,
				PotassiumSeverity:    SeverityLevel(sensor.PotassiumSeverity),
			},
			HydroponicManagerRelay: HydroponicManagerRelay{
				IsOn:                relay.IsOn,
				NextToggleInSeconds: relay.NextToggleInSeconds,
			},
		}
	default:
		fmt.Printf("Unsupported Hydroponic Manager message version: %d\n", payloadVersion)
		return nil
	}

	return &HydroponicManagerSensorDataResponse{
		PayloadVersion: HydroponicManagerMessageV1,
	}
}
