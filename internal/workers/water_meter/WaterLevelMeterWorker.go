package workers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/database"
)

const WaterLevelMeterTopicID int = 1

type WaterLevelMeterData struct {
	AverageWaterLevelCm string    `json:"waterLevelCm"`
	WaterLevelHistory   []float32 `json:"waterLevelHistory"`
	UptimeSeconds       int       `json:"uptimeSeconds"`
}

type WaterLevelMeterMessage struct {
	ClientId string              `json:"clientId"`
	FuseId   string              `json:"fuseId"`
	Data     WaterLevelMeterData `json:"data"`
}

type WaterLevelMeterListener struct {
	db *database.Database
}

func NewWaterLevelMeterListener(db *database.Database) *WaterLevelMeterListener {
	return &WaterLevelMeterListener{db: db}
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

	var message *WaterLevelMeterMessage = nil

	switch version {
	case 1:
		message, err = ParseWaterLevelMeterV1(parts)
		if err != nil {
			fmt.Printf("Failed to parse water meter message: %v\n", err)
			return
		}
	default:
		fmt.Printf("Unsupported water meter message version: %d\n", version)
		return
	}

	fmt.Printf("Parsed water meter message: %+v\n", message)

	dr := wm.db.DeviceRepository()
	sr := wm.db.SensorRepository()

	fuseID := message.ClientId
	ctx := context.Background()

	device, err := dr.GetDeviceByFuseID(ctx, fuseID)
	if err != nil {
		fmt.Printf("Device with client ID %s not found, creating new device\n", fuseID)
		device, err = dr.InsertDevice(context.Background(), fuseID, "Water Level Meter", "Water level meter device")
		if err != nil {
			fmt.Printf("Failed to insert device: %v\n", err)
			return
		}
	}

	err = sr.InsertSensorData(context.Background(), device.ID, WaterLevelMeterTopicID, "TODO", version)

	if err != nil {
		fmt.Printf("Failed to insert sensor data for device %s: %v\n", fuseID, err)
		return
	}

}

// 1;water-tank;39620398887400;waterLevelHistory:116.04:116.14:117.11:114.30:117.01:115.27:114.40:115.75:115.85:116.91;averageWaterLevelCm:115.88;uptime:60
func ParseWaterLevelMeterV1(parts []string) (*WaterLevelMeterMessage, error) {

	var message WaterLevelMeterMessage
	message.ClientId = parts[1]
	message.FuseId = parts[2]

	for i := 3; i < len(parts); i++ {
		var values = strings.Split(parts[i], ":")
		var key = values[0]
		if key == "waterLevelHistory" {
			messages := values[1:]
			var history []float32
			for _, v := range messages {
				value, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid water level history value: %s", v)
				}
				history = append(history, float32(value))
			}
			message.Data.WaterLevelHistory = history
			continue
		}
		if key == "averageWaterLevelCm" {
			value, err := strconv.ParseFloat(values[1], 32)
			if err != nil {
				return nil, fmt.Errorf("invalid average water level value: %s", values[1])
			}
			message.Data.AverageWaterLevelCm = fmt.Sprintf("%.2f", value)
			continue
		}

		if key == "uptime" {
			value, err := strconv.Atoi(values[1])
			if err != nil {
				return nil, fmt.Errorf("invalid uptime value: %s", values[1])
			}
			message.Data.UptimeSeconds = value
			continue
		}
	}

	return &message, nil
}
