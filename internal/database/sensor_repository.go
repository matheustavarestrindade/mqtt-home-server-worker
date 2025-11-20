package database

import (
	"context"
	"fmt"
	"time"
)

type SensorData struct {
	ID             int       `json:"id"`
	DeviceID       int       `json:"device_id"`
	TopicID        int       `json:"topic_id"`
	Payload        string    `json:"payload"`
	PayloadVersion int       `json:"payload_version"`
	CreatedAt      time.Time `json:"created_at"`
}

type SensorRepository struct {
	db *Database
}

func newSensorRepository(db *Database) *SensorRepository {
	return &SensorRepository{db: db}
}

func (r *SensorRepository) InsertSensorData(ctx context.Context, deviceID int, topicId int, payload string, payloadVersion int) error {
	_, err := r.db.pool.Exec(ctx, `
		WITH inserted_data AS (
			INSERT INTO sensor_data
				(device_id, topic_id, payload, payload_version)
			VALUES
				($1, $2, $3, $4)
		)
		UPDATE devices
		SET last_seen = NOW()
		WHERE id = $1;
	`, deviceID, topicId, payload, payloadVersion)
	if err != nil {
		return fmt.Errorf("failed to insert sensor data: %w", err)
	}

	return nil
}

func (r *SensorRepository) GetSensorDataByDeviceID(ctx context.Context, deviceID int) ([]SensorData, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, device_id, topic_id, payload, created_at 
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
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, device_id, topic_id, payload, payload_version, created_at 
		FROM sensor_data 
		WHERE device_id = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at ASC
	`, deviceID, startTime, endTime)

	if err != nil {
		return nil, fmt.Errorf("failed to query sensor data: %w", err)
	}

	defer rows.Close()

	var sensorData []SensorData
	for rows.Next() {
		var data SensorData
		err := rows.Scan(&data.ID, &data.DeviceID, &data.TopicID, &data.Payload, &data.PayloadVersion, &data.CreatedAt)
		sensorData = append(sensorData, data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sensor data: %w", err)
		}
	}

	return sensorData, nil
}
