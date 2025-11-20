package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/database"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/http"
	"github.com/matheustavarestrindade/mqtt-home-server-worker/internal/services"
	hydroponic_manager_worker "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/hydroponic_manager"
	water_meter_worker "github.com/matheustavarestrindade/mqtt-home-server-worker/internal/workers/water_meter"
)

type Config struct {
	DatabaseUrl string `env:"DATABASE_URL"`
	ClientId    string `env:"MQTT_CLIENT_ID"`
}
type Instance struct {
	Config     Config
	Database   *database.Database
	MQTTClient *services.MQTTClient
	HTTPServer *http.Server
}

var instance *Instance

func GetInstance() *Instance {
	if instance == nil {
		config := loadEnv()
		db := database.New()

		instance = &Instance{
			Config:   config,
			Database: db,
			MQTTClient: services.NewMQTTClient(services.MQTTConfig{
				BrokerUrl:         "mqtts://mosquitto.trindademedia.dev:8883",
				ClientId:          config.ClientId, 
				CaFilePath:        "./certs/ca.crt",
				ClientCrtFilePath: "./certs/client.crt",
				ClientKeyFilePath: "./certs/client.key",
			}),
			HTTPServer: http.NewServer(3000, db),
		}
	}

	return instance
}

func main() {
	fmt.Println("[MQTT Worker] Starting Worker")

	instance = GetInstance()

	err := instance.Database.Connect(instance.Config.DatabaseUrl)
	AssertOrExit(err, "Failed to connect to database with URL: %s", instance.Config.DatabaseUrl)

	err = instance.Database.RunMigrations()
	AssertOrExit(err, "Failed to run database migrations")

	err = instance.MQTTClient.Start()
	AssertOrExit(err, "Failed to start MQTT worker")

	hydroponic_manager_worker.NewHydroponicManagerListener(instance.Database, instance.MQTTClient)
	water_meter_worker.NewWaterLevelMeterListener(instance.Database, instance.MQTTClient)

	for instance.MQTTClient.IsRunning() {
		time.Sleep(200 * time.Millisecond)
	}
}

func AssertOrExit(err error, message string, vars ...any) {
	if err != nil {
		if len(vars) > 0 {
			fmt.Printf(message+"\n", vars...)
		} else {
			fmt.Println(message)
		}
		fmt.Println("Error:", err)
		panic(err)
	}
}

func loadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, reading configuration from environment variables")
	}

	return Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		ClientId:    os.Getenv("MQTT_CLIENT_ID"),
	}
}
