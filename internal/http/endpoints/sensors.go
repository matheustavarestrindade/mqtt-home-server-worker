package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/database"
)

type SensorEndpoints struct {
	db *database.Database
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
	fuseId, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(rw, "Invalid fuse ID", http.StatusBadRequest)
		return
	}
	startTimeUnix, err := strconv.ParseInt(r.URL.Query().Get("start"), 10, 64)
	if err != nil {
		http.Error(rw, "Invalid start time", http.StatusBadRequest)
		return
	}
	endTimeUnix, err := strconv.ParseInt(r.URL.Query().Get("end"), 10, 64)
	if err != nil {
		http.Error(rw, "Invalid end time", http.StatusBadRequest)
		return
	}

	startTime := time.Unix(startTimeUnix, 0)
	endTime := time.Unix(endTimeUnix, 0)

	// Start time must be before end time
	if !startTime.Before(endTime) {
		http.Error(rw, "Start time must be before end time", http.StatusBadRequest)
		return
	}

	sensor, err := se.db.DeviceRepository().GetDeviceByFuseID(r.Context(), strconv.FormatInt(fuseId, 10))
	if err != nil {
		http.Error(rw, "Failed to get sensor", http.StatusInternalServerError)
		return
	}
	sensorData, err := se.db.SensorRepository().GetSensorDataByDeviceIDWithTimestamp(r.Context(), sensor.ID, startTime, endTime)
	if err != nil {
		http.Error(rw, "Failed to get sensor data", http.StatusInternalServerError)
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
