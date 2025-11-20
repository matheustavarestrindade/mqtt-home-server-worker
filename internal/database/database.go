package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool             *pgxpool.Pool
	sensorRepository *SensorRepository
	deviceRepository *DeviceRepository
}

func New() *Database {
	return &Database{}
}

func (db *Database) Connect(databaseUrl string) error {
	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return err
	}
	db.pool = pool
	return nil
}

func (db *Database) SensorRepository() *SensorRepository {
	if db.sensorRepository == nil {
		db.sensorRepository = newSensorRepository(db)
	}
	return db.sensorRepository
}

func (db *Database) DeviceRepository() *DeviceRepository {
	if db.deviceRepository == nil {
		db.deviceRepository = newDeviceRepository(db)
	}
	return db.deviceRepository
}

func (db *Database) Close() error {
	db.pool.Close()
	return nil
}
