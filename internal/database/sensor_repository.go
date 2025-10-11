package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type SensorData struct {
	ID             int    `json:"id"`
	DeviceID       int    `json:"device_id"`
	TopicID        int    `json:"topic"`
	Payload        string `json:"payload"`
	PayloadVersion int    `json:"payload_version"`
	CreatedAt      string `json:"created_at"`
}

type SensorRepository struct {
	db *Database
}

func newSensorRepository(db *Database) *SensorRepository {
	return &SensorRepository{db: db}
}

func (r *SensorRepository) InsertSensorData(ctx context.Context, deviceID int, topic int, payload string, payloadVersion int) error {
	jsonData, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = r.db.conn.Exec(ctx, `
		INSERT INTO sensor_data 
			(device_id, topic_id, payload, payload_version)
		VALUES 
			($1, $2, $3, $4)
	`, deviceID, topic, jsonData, payloadVersion)
	if err != nil {
		return fmt.Errorf("failed to insert sensor data: %w", err)
	}

	return nil
}

func (r *SensorRepository) GetSensorDataByDeviceID(ctx context.Context, deviceID int) ([]SensorData, error) {
	rows, err := r.db.conn.Query(ctx, `
		SELECT id, device_id, topic, payload, created_at 
		FROM sensor_data 
		WHERE device_id = $1
	`, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sensor data: %w", err)
	}
	defer rows.Close()

	var sensorData []SensorData
	for rows.Next() {
		var data SensorData
		if err := rows.Scan(&data.ID, &data.DeviceID, &data.TopicID, &data.Payload, &data.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan sensor data: %w", err)
		}
		sensorData = append(sensorData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return sensorData, nil
}

func (r *SensorRepository) GetSensorDataByDeviceIDWithTimestamp(ctx context.Context, deviceID int, startTime time.Time, endTime time.Time) ([]SensorData, error) {
	rows, err := r.db.conn.Query(ctx, `
		SELECT id, device_id, topic, payload, created_at 
		FROM sensor_data 
		WHERE device_id = $1 AND created_at >= $2 AND created_at <= $3
	`, deviceID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query sensor data: %w", err)
	}
	defer rows.Close()

	var sensorData []SensorData
	for rows.Next() {
		var data SensorData
		if err := rows.Scan(&data.ID, &data.DeviceID, &data.TopicID, &data.Payload, &data.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan sensor data: %w", err)
		}
		sensorData = append(sensorData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return sensorData, nil
}
