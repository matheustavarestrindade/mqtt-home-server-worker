package database

import (
	"context"
	"fmt"
	"time"
)

type Device struct {
	ID             int       `json:"id"`
	FuseID         string    `json:"fuse_id"`
	Name           string    `json:"name"`
	WifiStrength   int       `json:"wifi_strength"`
	BatteryPercent int       `json:"battery_percent"`
	Location       string    `json:"location"`
	Type           string    `json:"type"`
	Description    string    `json:"description"`
	LastSeen       time.Time `json:"last_seen"`
	CreatedAt      time.Time `json:"created_at"`
}

type DeviceRepository struct {
	db *Database
}

func newDeviceRepository(db *Database) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) InsertDevice(ctx context.Context, fuseID, name, description, location, deviceType string, wifiStrength, batteryPercent int) (*Device, error) {
	var device Device
	err := r.db.pool.QueryRow(ctx, `
		INSERT INTO devices 
			(fuseId, name, description, location, type, wifi_strength, battery_percent, last_seen)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING 
			id, fuseId, name, description, created_at, location, type, wifi_strength, battery_percent, last_seen
	`, fuseID, name, description, location, deviceType, wifiStrength, batteryPercent).Scan(
		&device.ID,
		&device.FuseID,
		&device.Name,
		&device.Description,
		&device.CreatedAt,
		&device.Location,
		&device.Type,
		&device.WifiStrength,
		&device.BatteryPercent,
		&device.LastSeen,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert device: %w", err)
	}

	return &device, nil
}

func (r *DeviceRepository) GetDeviceByFuseID(ctx context.Context, fuseID string) (*Device, error) {
	var device Device
	err := r.db.pool.QueryRow(ctx, `
		SELECT 
			id, fuseId, name, description, created_at, location, type, wifi_strength, battery_percent
		FROM 
			devices 
		WHERE 
			fuseId = $1
	`, fuseID).Scan(&device.ID, &device.FuseID, &device.Name, &device.Description, &device.CreatedAt, &device.Location, &device.Type, &device.WifiStrength, &device.BatteryPercent)

	if err != nil {
		return nil, fmt.Errorf("failed to query device: %w", err)
	}

	return &device, nil
}

func (r *DeviceRepository) GetDevicesByFuseID(ctx context.Context, fuseIds []string) ([]Device, error) {
	if len(fuseIds) == 0 {
		return []Device{}, nil
	}

	var devices []Device

	rows, err := r.db.pool.Query(ctx, `
		SELECT
			id, fuseId, name, description, created_at, location, type, wifi_strength, battery_percent, last_seen
		FROM
			devices
		WHERE
			fuseId = any($1)
	`, fuseIds)

	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var device Device
		err := rows.Scan(&device.ID,
			&device.FuseID,
			&device.Name,
			&device.Description,
			&device.CreatedAt,
			&device.Location,
			&device.Type,
			&device.WifiStrength,
			&device.BatteryPercent,
			&device.LastSeen)

		fmt.Println(device)

		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (dr *DeviceRepository) CreateAndGetDeviceIfDoesNotExist(fuseId string, name, description, location, deviceType string, wifiStrength, batteryPercent int) (*Device, error) {
	ctx := context.Background()

	device, err := dr.GetDeviceByFuseID(ctx, fuseId)
	if err != nil {
		fmt.Printf("Device with client ID %s not found, creating new device\n", fuseId)
		device, err = dr.InsertDevice(ctx, fuseId, name, description, location, deviceType, wifiStrength, batteryPercent)
		if err != nil {
			return nil, fmt.Errorf("failed to insert device: %v", err)
		}
	}

	return device, nil
}
