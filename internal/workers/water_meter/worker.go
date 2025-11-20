package water_meter_worker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/database"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/services"
	wm_payload_v1 "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/water_meter/payloads/water_meter_payload_v1"
)

const WaterLevelMeterTopicID = 1

const DeviceName = "Water Level Meter"
const DeviceDescription = "Water Level Meter Device"
const DeviceType = "water-level-meter"

// List of suported Paylaods versions
const (
	WaterMeterMessageV1 = 1
)

type WaterLevelMeterListener struct {
	db     *database.Database
	client *services.MQTTClient
}

func NewWaterLevelMeterListener(db *database.Database, client *services.MQTTClient) *WaterLevelMeterListener {
	hm := &WaterLevelMeterListener{
		db:     db,
		client: client,
	}

	client.Subscribe("water-meter/sensors", hm.Handler)
	return hm
}

func (wm *WaterLevelMeterListener) Handler(msg mqtt.Message) {

	fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))

	var payload string = string(msg.Payload())
	var parts = strings.Split(payload, ";")

	if len(parts) < 1 {
		fmt.Printf("Invalid water meter message format: %s\n", payload)
		return
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("Invalid water meter message version: %s\n", parts[0])
		return
	}

	var compressedData string
	var fuseID string
	var payloadVersion int = version

	switch version {
	case WaterMeterMessageV1:
		message, err := wm_payload_v1.ParsePayload(parts)
		if err != nil {
			fmt.Printf("Failed to parse water meter message: %v\n", err)
			return
		}
		fuseID = message.FuseId

		fmt.Printf("Parsed water meter v1 message: %+v\n", message)
		compressedData, err = wm_payload_v1.CompressDataToDatabase(message.Data)
		if err != nil {
			fmt.Printf("Failed to compress water meter v1 data: %v\n", err)
			return
		}
	default:
		fmt.Printf("Unsupported water meter message version: %d\n", version)
		return
	}

	sr := wm.db.SensorRepository()
	dr := wm.db.DeviceRepository()

	device, err := dr.CreateAndGetDeviceIfDoesNotExist(fuseID, DeviceName, DeviceDescription, "Unknown", DeviceType, 0, 0)
	if err != nil {
		fmt.Printf("Failed to insert device: %v\n", err)
		return
	}

	err = sr.InsertSensorData(context.Background(), device.ID, WaterLevelMeterTopicID, compressedData, payloadVersion)
	if err != nil {
		fmt.Printf("Failed to insert sensor data for device %s: %v\n", fuseID, err)
		return
	}
}
