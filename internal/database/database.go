package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	conn             pgx.Conn
	sensorRepository *SensorRepository
	deviceRepository *DeviceRepository
}

func New() *Database {
	return &Database{}
}

func (db *Database) Connect(databaseUrl string) error {
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		return err
	}
	db.conn = *conn
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
	if db.conn.IsClosed() {
		return nil
	}
	return db.conn.Close(context.Background())
}
