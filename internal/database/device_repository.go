package database

import (
	"context"
	"fmt"
	"time"
)

type Device struct {
	ID          int       `json:"id"`
	FuseID      string    `json:"fuse_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type DeviceRepository struct {
	db *Database
}

func newDeviceRepository(db *Database) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) InsertDevice(ctx context.Context, fuseID, name, description string) (*Device, error) {
	var device Device
	err := r.db.conn.QueryRow(ctx, `
		INSERT INTO devices 
			(fuseId, name, description)
		VALUES 
			($1, $2, $3)
		RETURNING 
			id, fuseId, name, description, created_at
	`, fuseID, name, description).Scan(
		&device.ID,
		&device.FuseID,
		&device.Name,
		&device.Description,
		&device.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert device: %w", err)
	}

	return &device, nil
}

func (r *DeviceRepository) GetDeviceByFuseID(ctx context.Context, fuseID string) (*Device, error) {
	var device Device
	err := r.db.conn.QueryRow(ctx, `
		SELECT 
			id, fuseId, name, description, created_at 
		FROM 
			devices 
		WHERE 
			fuseId = $1
	`, fuseID).Scan(&device.ID, &device.FuseID, &device.Name, &device.Description, &device.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to query device: %w", err)
	}

	return &device, nil
}

func (r *DeviceRepository) GetDevicesByFuseID(ctx context.Context, clientIDs []string) ([]Device, error) {
	if len(clientIDs) == 0 {
		return []Device{}, nil
	}

	var devices []Device
	err := r.db.conn.QueryRow(ctx, `
		SELECT
			id, fuseId, name, description, created_at
		FROM
			devices
		WHERE
			fuseId in ($1)
	`, clientIDs).Scan(&devices)

	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}

	return devices, nil
}
