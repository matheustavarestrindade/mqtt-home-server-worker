package hydroponic_manager_payload_v1

import (
	"fmt"
	"strconv"
	"strings"
)

// Commands List
// Must be appended with: payloadVersion;clientId;ESP.fuseMac;action:parameters
// List of actions:parameters bellow
//
// toggle_relay
// set_nitrogen_thresholds:min,max
// set_phosphorus_thresholds:min,max
// set_potassium_thresholds:min,max
// set_ph_thresholds:min,max
// set_conductivity_thresholds:min,max
// set_water_level_thresholds:min,max
// set_percentage_thresholds:warn,critical
// set_wifi_credentials:ssid,password
// set_automatic_configuration:true/false
// set_crops_ids:id1,id2,id3
// restart_device
//
// Payload example:
//
// String payload = String(MQTT_MESSAGE_VERSION) + ";";
// payload += String(clientId) + ";";
// payload += String(ESP.getEfuseMac()) + ";";
// payload += "relay:" +
//     String(rs.isOn) + ":" +
//     String(rs.nextSwitchTimeInSeconds) + ";";
// payload += "sensor:" +
//     String(sd.moisture) + ":" + static_cast<int>(sd.moistureSeverity) + ":" +
//     String(sd.temperature) + ":" + static_cast<int>(sd.temperatureSeverity) + ":" +
//     String(sd.conductivity) + ":" + static_cast<int>(sd.conductivitySeverity) + ":" +
//     String(sd.ph) + ":" + static_cast<int>(sd.phSeverity) + ":" +
//     String(sd.nitrogen) + ":" + static_cast<int>(sd.nitrogenSeverity) + ":" +
//     String(sd.phosphorus) + ":" + static_cast<int>(sd.phosphorusSeverity) + ":" +
//     String(sd.potassium) + ":" + static_cast<int>(sd.potassiumSeverity) + ";";
// payload += "water_level:" + String(wl.levelCm);

type Command string

const (
	CommandToggleRelay               Command = "toggle_relay"
	CommandSetNitrogenThresholds     Command = "set_nitrogen_thresholds"
	CommandSetPhosphorusThresholds   Command = "set_phosphorus_thresholds"
	CommandSetPotassiumThresholds    Command = "set_potassium_thresholds"
	CommandSetPhThresholds           Command = "set_ph_thresholds"
	CommandSetConductivityThresholds Command = "set_conductivity_thresholds"
	CommandSetWaterLevelThresholds   Command = "set_water_level_thresholds"
	CommandSetPercentageThresholds   Command = "set_percentage_thresholds"
	CommandSetWifiCredentials        Command = "set_wifi_credentials"
	CommandSetAutomaticConfiguration Command = "set_automatic_configuration"
	CommandSetCropsIds               Command = "set_crops_ids"
	CommandRestartDevice             Command = "restart_device"
)

type Payload struct {
	Version  int    `json:"version"`
	ClientId string `json:"clientId"`
	FuseId   string `json:"espFuseId"`
	Data     Data   `json:"data"`
}

type Data struct {
	Sensors SensorData             `json:"sensors"`
	Relay   HydroponicManagerRelay `json:"relay"`
}

type SeverityLevel int

const MAX_COMPRESSED_PAYLOAD_LENGTH = 128

const (
	NORMAL SeverityLevel = iota
	WARNING
	CRITICAL
)

type SensorData struct {
	Temperature          float32       `json:"temperature"`
	TemperaturaSeverity  SeverityLevel `json:"temperatureSeverity"`
	Moisture             float32       `json:"moisture"`
	MoistureSeverity     SeverityLevel `json:"moistureSeverity"`
	Ph                   float32       `json:"ph"`
	PhSeverity           SeverityLevel `json:"phSeverity"`
	Conductivity         int           `json:"conductivity"`
	ConductivitySeverity SeverityLevel `json:"conductivitySeverity"`
	Nitrogen             int           `json:"nitrogen"`
	NitrogenSeverity     SeverityLevel `json:"nitrogenSeverity"`
	Phosphorus           int           `json:"phosphorus"`
	PhosphorusSeverity   SeverityLevel `json:"phosphorusSeverity"`
	Potassium            int           `json:"potassium"`
	PotassiumSeverity    SeverityLevel `json:"potassiumSeverity"`
}

type HydroponicManagerRelay struct {
	IsOn                bool `json:"isOn"`
	NextToggleInSeconds int  `json:"nextToggleInSeconds"`
}

func ParsePayload(parts []string) (*Payload, error) {
	var message Payload

	message.ClientId = parts[1]
	message.FuseId = parts[2]

	for i := 3; i < len(parts); i++ {
		var values = strings.Split(parts[i], ":")
		var key = values[0]

		if key == "sensor" {
			moisture, err := strconv.ParseFloat(values[1], 32)
			if err != nil {
				return nil, fmt.Errorf("invalid sensor moisture value: %s", values[1])
			}
			moistureSeverity, err := strconv.Atoi(values[2])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor moisture severity value: %s", values[2])
			}
			temperature, err := strconv.ParseFloat(values[3], 32)
			if err != nil {
				return nil, fmt.Errorf("invalid sensor temperature value: %s", values[3])
			}
			temperatureSeverity, err := strconv.Atoi(values[4])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor temperature severity value: %s", values[4])
			}
			conductivity, err := strconv.Atoi(values[5])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor conductivity value: %s", values[5])
			}
			conductivitySeverity, err := strconv.Atoi(values[6])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor conductivity severity value: %s", values[6])
			}
			ph, err := strconv.ParseFloat(values[7], 32)
			if err != nil {
				return nil, fmt.Errorf("invalid sensor ph value: %s", values[7])
			}
			phSeverity, err := strconv.Atoi(values[8])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor ph severity value: %s", values[8])
			}
			nitrogen, err := strconv.Atoi(values[9])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor nitrogen value: %s", values[9])
			}
			nitrogenSeverity, err := strconv.Atoi(values[10])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor nitrogen severity value: %s", values[10])
			}
			phosphorus, err := strconv.Atoi(values[11])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor phosphorus value: %s", values[11])
			}
			phosphorusSeverity, err := strconv.Atoi(values[12])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor phosphorus severity value: %s", values[12])
			}
			potassium, err := strconv.Atoi(values[13])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor potassium value: %s", values[13])
			}
			potassiumSeverity, err := strconv.Atoi(values[14])
			if err != nil {
				return nil, fmt.Errorf("invalid sensor potassium severity value: %s", values[14])
			}

			message.Data.Sensors = SensorData{
				Temperature:          float32(temperature),
				TemperaturaSeverity:  SeverityLevel(temperatureSeverity),
				Moisture:             float32(moisture),
				MoistureSeverity:     SeverityLevel(moistureSeverity),
				Ph:                   float32(ph),
				PhSeverity:           SeverityLevel(phSeverity),
				Conductivity:         conductivity,
				ConductivitySeverity: SeverityLevel(conductivitySeverity),
				Nitrogen:             nitrogen,
				NitrogenSeverity:     SeverityLevel(nitrogenSeverity),
				Phosphorus:           phosphorus,
				PhosphorusSeverity:   SeverityLevel(phosphorusSeverity),
				Potassium:            potassium,
				PotassiumSeverity:    SeverityLevel(potassiumSeverity),
			}
			continue
		}

		if key == "relay" {
			isOn, err := strconv.ParseBool(values[1])
			if err != nil {
				return nil, fmt.Errorf("invalid relay isOn value: %s", values[1])
			}
			nextToggleInSeconds, err := strconv.Atoi(values[2])
			if err != nil {
				return nil, fmt.Errorf("invalid relay nextToggleInSeconds value: %s", values[2])
			}
			message.Data.Relay = HydroponicManagerRelay{
				IsOn:                isOn,
				NextToggleInSeconds: nextToggleInSeconds,
			}
			continue
		}
	}

	return &message, nil
}

func CreateCommand(clientId string, fuseId string, command Command, args []string) string {
	return fmt.Sprintf("1;%s;%s;%s:%s", clientId, fuseId, command, strings.Join(args, ","))
}

func CompressDataToDatabase(data Data) (string, error) {
	compressedData := fmt.Sprintf("T:%.2f:%d;M:%.2f:%d;pH:%.2f:%d;C:%d:%d;N:%d:%d;P:%d:%d;K:%d:%d;R:%t:%d",
		data.Sensors.Temperature, data.Sensors.TemperaturaSeverity,
		data.Sensors.Moisture, data.Sensors.MoistureSeverity,
		data.Sensors.Ph, data.Sensors.PhSeverity,
		data.Sensors.Conductivity, data.Sensors.ConductivitySeverity,
		data.Sensors.Nitrogen, data.Sensors.NitrogenSeverity,
		data.Sensors.Phosphorus, data.Sensors.PhosphorusSeverity,
		data.Sensors.Potassium, data.Sensors.PotassiumSeverity,
		data.Relay.IsOn, data.Relay.NextToggleInSeconds,
	)

	if len(compressedData) > MAX_COMPRESSED_PAYLOAD_LENGTH {
		fmt.Printf("Warning: Compressed data length %d exceeds maximum of %d characters\n", len(compressedData), MAX_COMPRESSED_PAYLOAD_LENGTH)
		return "", fmt.Errorf("compressed data exceeds maximum length of %d characters", MAX_COMPRESSED_PAYLOAD_LENGTH)
	}
	return compressedData, nil
}

func DecompressDataFromDatabase(compressedData string) (Data, error) {
	var data Data

	for part := range strings.SplitSeq(compressedData, ";") {
		values := strings.Split(part, ":")
		key := values[0]

		switch key {
		case "T":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid temperature data format")
			}
			temp, err := strconv.ParseFloat(values[1], 32)
			if err != nil {
				return data, fmt.Errorf("invalid temperature value: %s", values[1])
			}
			severity, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid temperature severity value: %s", values[2])
			}
			data.Sensors.Temperature = float32(temp)
			data.Sensors.TemperaturaSeverity = SeverityLevel(severity)
		case "M":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid moisture data format")
			}
			moisture, err := strconv.ParseFloat(values[1], 32)
			if err != nil {
				return data, fmt.Errorf("invalid moisture value: %s", values[1])
			}
			severity, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid moisture severity value: %s", values[2])
			}
			data.Sensors.Moisture = float32(moisture)
			data.Sensors.MoistureSeverity = SeverityLevel(severity)
		case "pH":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid pH data format")
			}
			ph, err := strconv.ParseFloat(values[1], 32)
			if err != nil {
				return data, fmt.Errorf("invalid pH value: %s", values[1])
			}
			severity, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid pH severity value: %s", values[2])
			}
			data.Sensors.Ph = float32(ph)
			data.Sensors.PhSeverity = SeverityLevel(severity)
		case "C":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid conductivity data format")
			}
			conductivity, err := strconv.Atoi(values[1])
			if err != nil {
				return data, fmt.Errorf("invalid conductivity value: %s", values[1])
			}
			severity, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid conductivity severity value: %s", values[2])
			}
			data.Sensors.Conductivity = conductivity
			data.Sensors.ConductivitySeverity = SeverityLevel(severity)
		case "N":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid nitrogen data format")
			}
			nitrogen, err := strconv.Atoi(values[1])
			if err != nil {
				return data, fmt.Errorf("invalid nitrogen value: %s", values[1])
			}
			severity, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid nitrogen severity value: %s", values[2])
			}
			data.Sensors.Nitrogen = nitrogen
			data.Sensors.NitrogenSeverity = SeverityLevel(severity)
		case "P":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid phosphorus data format")
			}
			phosphorus, err := strconv.Atoi(values[1])
			if err != nil {
				return data, fmt.Errorf("invalid phosphorus value: %s", values[1])
			}
			severity, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid phosphorus severity value: %s", values[2])
			}
			data.Sensors.Phosphorus = phosphorus
			data.Sensors.PhosphorusSeverity = SeverityLevel(severity)
		case "K":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid potassium data format")
			}
			potassium, err := strconv.Atoi(values[1])
			if err != nil {
				return data, fmt.Errorf("invalid potassium value: %s", values[1])
			}
			severity, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid potassium severity value: %s", values[2])
			}
			data.Sensors.Potassium = potassium
			data.Sensors.PotassiumSeverity = SeverityLevel(severity)
		case "R":
			if len(values) != 3 {
				return data, fmt.Errorf("invalid relay data format")
			}
			isOn, err := strconv.ParseBool(values[1])
			if err != nil {
				return data, fmt.Errorf("invalid relay isOn value: %s", values[1])
			}
			nextToggleInSeconds, err := strconv.Atoi(values[2])
			if err != nil {
				return data, fmt.Errorf("invalid relay nextToggleInSeconds value: %s", values[2])
			}
			data.Relay.IsOn = isOn
			data.Relay.NextToggleInSeconds = nextToggleInSeconds
		default:
			fmt.Printf("Warning: Unknown data key: %s\n", key)
			fmt.Printf("Values for unknown key: %v\n", values)
			return data, fmt.Errorf("unknown data key: %s", key)
		}
	}

	return data, nil
}
