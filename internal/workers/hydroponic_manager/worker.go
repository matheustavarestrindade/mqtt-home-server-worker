package hydroponic_manager_worker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/database"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/services"
	hm_payload_v1 "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/hydroponic_manager/payloads/hydroponic_manager_payload_v1"
)

const HydroponicManagerTopicID = 0
const DeviceName = "Hydroponic Manager"
const DeviceType = "hydroponic-manager"
const DeviceDescription = "Hydroponic Manager Device"

// List of suported Paylaods versions
const (
	HydroponicManagerMessageV1 = 1
)

type HydroponicManagerWorker struct {
	db     *database.Database
	client *services.MQTTClient
}

func NewHydroponicManagerListener(db *database.Database, client *services.MQTTClient) *HydroponicManagerWorker {
	hm := &HydroponicManagerWorker{
		db:     db,
		client: client,
	}

	client.Subscribe("hydroponic-manager/sensors", hm.Handler)
	return hm
}

func (hm *HydroponicManagerWorker) Handler(msg mqtt.Message) {
	fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))

	var payload string = string(msg.Payload())
	var parts = strings.Split(payload, ";")

	if len(parts) < 1 {
		fmt.Printf("Invalid Hydroponic Manager message format: %s\n", payload)
		return
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("Invalid Hydroponic Manager message version: %s\n", parts[0])
		return
	}

	var compressedData string
	var fuseID string
	var payloadVersion int = version

	switch version {
	case HydroponicManagerMessageV1:
		message, err := hm_payload_v1.ParsePayload(parts)
		if err != nil {
			fmt.Printf("Failed to parse Hydroponic Manager message: %v\n", err)
			return
		}
		fuseID = message.FuseId
		compressedData, err = hm_payload_v1.CompressDataToDatabase(message.Data)
		if err != nil {
			fmt.Printf("Failed to compress Hydroponic Manager v1 data: %v\n", err)
			return
		}
		fmt.Printf("Parsed Hydroponic Manager v1 Payload: %+v\n", message)
	default:
		fmt.Printf("Unsupported water Hydroponic Manager meter message version: %d\n", version)
		return
	}

	sr := hm.db.SensorRepository()
	dr := hm.db.DeviceRepository()

	// TODO:Parse payload to get location and wifi/battery status
	device, err := dr.CreateAndGetDeviceIfDoesNotExist(fuseID, DeviceName, DeviceDescription, "Unknown", DeviceType, 0, 0)
	if err != nil {
		fmt.Printf("Failed to insert device: %v\n", err)
		return
	}

	err = sr.InsertSensorData(context.Background(), device.ID, HydroponicManagerTopicID, compressedData, payloadVersion)

	if err != nil {
		fmt.Printf("Failed to insert sensor data for device %s: %v\n", fuseID, err)
		return
	}

}
