package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/database"
	hydroponic_manager_worker "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/hydroponic_manager"
	water_meter_worker "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/water_meter"
)

type SensorEndpoints struct {
	db *database.Database
}

type SumData struct {
	Temperature          float32
	TemperaturaSeverity  int
	Moisture             float32
	MoistureSeverity     int
	Ph                   float32
	PhSeverity           int
	Conductivity         int
	ConductivitySeverity int
	Nitrogen             int
	NitrogenSeverity     int
	Phosphorus           int
	PhosphorusSeverity   int
	Potassium            int
	PotassiumSeverity    int
}

func NewSensorEndpoints(db *database.Database) *SensorEndpoints {
	return &SensorEndpoints{db: db}
}

func (se *SensorEndpoints) GetSensorsByID(rw http.ResponseWriter, r *http.Request) {
	fuseIds := strings.Split(r.URL.Query().Get("ids"), ",")
	if len(fuseIds) == 0 {
		http.Error(rw, "No fuse IDs provided", http.StatusBadRequest)
		return
	}

	sensors, err := se.db.DeviceRepository().GetDevicesByFuseID(r.Context(), fuseIds)
	if err != nil {
		http.Error(rw, "Failed to get sensors", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(sensors)
	if err != nil {
		http.Error(rw, "Failed to marshal sensors", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
}

func (se *SensorEndpoints) GetSensorDataByIDAndTimestamp(rw http.ResponseWriter, r *http.Request) {

	fuseId := r.URL.Query().Get("fuse_id")
	startTimeISO := r.URL.Query().Get("start")
	endTimeISO := r.URL.Query().Get("end")
	interval_ms, err := strconv.Atoi(r.URL.Query().Get("interval_ms"))

	startTime, err := time.Parse(time.RFC3339, startTimeISO)
	if err != nil {
		http.Error(rw, "Invalid start time format. Use ISO 8601 format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeISO)
	if err != nil {
		http.Error(rw, "Invalid end time format. Use ISO 8601 format", http.StatusBadRequest)
		return
	}

	// Start time must be before end time
	if !startTime.Before(endTime) {
		http.Error(rw, "Start time must be before end time", http.StatusBadRequest)
		return
	}

	device, err := se.db.DeviceRepository().GetDeviceByFuseID(r.Context(), fuseId)
	if err != nil {
		http.Error(rw, "Failed to get device by fuse ID", http.StatusInternalServerError)
		return
	}
	sensorDataCompressed, err := se.db.SensorRepository().GetSensorDataByDeviceIDWithTimestamp(r.Context(), device.ID, startTime, endTime)
	if err != nil {
		fmt.Println("Error fetching sensor data:", err)
		http.Error(rw, "Failed to get sensor data", http.StatusInternalServerError)
		return
	}
	if len(sensorDataCompressed) == 0 {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("[]"))
		return
	}

	var sensorData any

	if device.Type == "hydroponic-manager" {
		sensorData, err = se.handleHydroponicManagerSensorDataAggregation(sensorDataCompressed, interval_ms, startTime, endTime, device.ID)
		if err != nil {
			http.Error(rw, "Failed to aggregate sensor data", http.StatusInternalServerError)
			return
		}
	} else if device.Type == "water-level-meter" {
		sensorData, err = handleWaterManagerSensorDataAggregation(sensorDataCompressed, interval_ms, startTime, endTime, device.ID)
		if err != nil {
			http.Error(rw, "Failed to aggregate sensor data", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(rw, "Unsupported device type or interval", http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(sensorData)
	if err != nil {
		http.Error(rw, "Failed to marshal sensors", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
}

func handleWaterManagerSensorDataAggregation(sensorDataCompressed []database.SensorData, interval_ms int, startTime, endTime time.Time, deviceID int) ([]water_meter_worker.WaterLevelMeterSensorDataResponse, error) {

	sensorData := make([]water_meter_worker.WaterLevelMeterSensorDataResponse, 0)

	waterLevelSum := 0.0
	currentTime := startTime.Add(time.Duration(interval_ms) * time.Millisecond)
	count := 0

	for _, data := range sensorDataCompressed {
		dataConverted := water_meter_worker.ConvertCompressedPayloadToSensorDataResponse(data.PayloadVersion, data.Payload)
		if dataConverted == nil {
			fmt.Printf("Failed to convert compressed payload to sensor data response for device ID %d\n", deviceID)
			continue
		}
		if interval_ms == 0 {
			sensorData = append(sensorData, *dataConverted)
			continue
		}

		for data.CreatedAt.After(currentTime) && count == 0 {
			sensorData = append(sensorData, water_meter_worker.WaterLevelMeterSensorDataResponse{})
			currentTime = currentTime.Add(time.Duration(interval_ms) * time.Millisecond)
		}

		waterLevelSum += float64(dataConverted.AverageWaterLevelCm)
		count++

		if data.CreatedAt.Before(currentTime) {
			continue
		}
		finalData := &water_meter_worker.WaterLevelMeterSensorDataResponse{}
		finalData.AverageWaterLevelCm = float32(waterLevelSum / float64(count))

		count = 0
		currentTime = currentTime.Add(time.Duration(interval_ms) * time.Millisecond)
		waterLevelSum = 0.0
		sensorData = append(sensorData, *finalData)
	}
	return sensorData, nil
}

func (se *SensorEndpoints) handleHydroponicManagerSensorDataAggregation(sensorDataCompressed []database.SensorData, interval_ms int, startTime, endTime time.Time, deviceID int) ([]hydroponic_manager_worker.HydroponicManagerSensorDataResponse, error) {
	sensorData := make([]hydroponic_manager_worker.HydroponicManagerSensorDataResponse, 0)

	currentDataSum := SumData{}
	currentTime := startTime.Add(time.Duration(interval_ms) * time.Millisecond)
	count := 0

	for _, data := range sensorDataCompressed {
		dataConverted := hydroponic_manager_worker.ConvertCompressedPayloadToSensorDataResponse(data.PayloadVersion, data.Payload)
		if dataConverted == nil {
			fmt.Printf("Failed to convert compressed payload to sensor data response for device ID %d\n", deviceID)
			continue
		}

		if interval_ms == 0 {
			sensorData = append(sensorData, *dataConverted)
			continue
		}

		// Fill empty intervals with empty data
		for data.CreatedAt.After(currentTime) && count == 0 {
			sensorData = append(sensorData, hydroponic_manager_worker.HydroponicManagerSensorDataResponse{})
			currentTime = currentTime.Add(time.Duration(interval_ms) * time.Millisecond)
		}

		fmt.Printf("interval_ms: %d\n", interval_ms)
		fmt.Printf("data.CreatedAt: %s\n", data.CreatedAt.String())
		fmt.Printf("currentTime: %s\n", currentTime.String())
		fmt.Printf("dataConverted: %+v\n", dataConverted)

		timeDiff := currentTime.Sub(data.CreatedAt)
		fmt.Printf("timeDiff: %s\n", timeDiff.String())

		currentDataSum.Temperature += dataConverted.Temperature
		currentDataSum.TemperaturaSeverity += int(dataConverted.TemperaturaSeverity)
		currentDataSum.Moisture += dataConverted.Moisture
		currentDataSum.MoistureSeverity += int(dataConverted.MoistureSeverity)
		currentDataSum.Ph += dataConverted.Ph
		currentDataSum.PhSeverity += int(dataConverted.PhSeverity)
		currentDataSum.Conductivity += dataConverted.Conductivity
		currentDataSum.ConductivitySeverity += int(dataConverted.ConductivitySeverity)
		currentDataSum.Nitrogen += dataConverted.Nitrogen
		currentDataSum.NitrogenSeverity += int(dataConverted.NitrogenSeverity)
		currentDataSum.Phosphorus += dataConverted.Phosphorus
		currentDataSum.PhosphorusSeverity += int(dataConverted.PhosphorusSeverity)
		currentDataSum.Potassium += dataConverted.Potassium
		currentDataSum.PotassiumSeverity += int(dataConverted.PotassiumSeverity)
		count++

		if data.CreatedAt.Before(currentTime) {
			continue
		}

		finalData := &hydroponic_manager_worker.HydroponicManagerSensorDataResponse{}

		finalData.Temperature = currentDataSum.Temperature / float32(count)
		finalData.TemperaturaSeverity = hydroponic_manager_worker.CalculateSeverityLevel(currentDataSum.TemperaturaSeverity / count)
		finalData.Moisture = currentDataSum.Moisture / float32(count)
		finalData.MoistureSeverity = hydroponic_manager_worker.CalculateSeverityLevel(currentDataSum.MoistureSeverity / count)
		finalData.Ph = currentDataSum.Ph / float32(count)
		finalData.PhSeverity = hydroponic_manager_worker.CalculateSeverityLevel(currentDataSum.PhSeverity / count)
		finalData.Conductivity = currentDataSum.Conductivity / count
		finalData.ConductivitySeverity = hydroponic_manager_worker.CalculateSeverityLevel(currentDataSum.ConductivitySeverity / count)
		finalData.Nitrogen = currentDataSum.Nitrogen / count
		finalData.NitrogenSeverity = hydroponic_manager_worker.CalculateSeverityLevel(currentDataSum.NitrogenSeverity / count)
		finalData.Phosphorus = currentDataSum.Phosphorus / count
		finalData.PhosphorusSeverity = hydroponic_manager_worker.CalculateSeverityLevel(currentDataSum.PhosphorusSeverity / count)
		finalData.Potassium = currentDataSum.Potassium / count
		finalData.PotassiumSeverity = hydroponic_manager_worker.CalculateSeverityLevel(currentDataSum.PotassiumSeverity / count)

		finalData.IsOn = dataConverted.IsOn
		finalData.NextToggleInSeconds = dataConverted.NextToggleInSeconds

		count = 0
		currentTime = currentTime.Add(time.Duration(interval_ms) * time.Millisecond)

		fmt.Printf("Final data: %+v\n", finalData)
		fmt.Printf("Current time: %s\n", currentTime.String())
		fmt.Printf("currentDataSum: %+v\n", currentDataSum)

		sensorData = append(sensorData, *finalData)
		currentDataSum = SumData{}
	}

	return sensorData, nil
}
