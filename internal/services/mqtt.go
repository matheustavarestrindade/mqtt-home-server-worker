package services

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConfig struct {
	CaFilePath        string
	ClientCrtFilePath string
	ClientKeyFilePath string
	ClientId          string
	BrokerUrl         string
}

type MqttMessageHandler func(msg mqtt.Message)

type MQTTClient struct {
	config   MQTTConfig
	client   mqtt.Client
	handlers map[string]MqttMessageHandler
}

func NewMQTTClient(config MQTTConfig) *MQTTClient {
	return &MQTTClient{
		config:   config,
		handlers: make(map[string]MqttMessageHandler),
	}
}

func (worker *MQTTClient) IsRunning() bool {
	return worker.client.IsConnected()
}

func (worker *MQTTClient) Start() error {
	config := worker.config

	caCert, err := os.ReadFile(config.CaFilePath)
	if err != nil {
		return fmt.Errorf("Error reading CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(config.ClientCrtFilePath, config.ClientKeyFilePath)
	if err != nil {
		return fmt.Errorf("Error reading client certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{clientCert},
		InsecureSkipVerify: true, // TODO: Need to fix SAN field on certificate to make this production ready
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.BrokerUrl)
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(config.ClientId)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	worker.client = mqtt.NewClient(opts)
	if token := worker.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("Error connecting to MQTT broker: %v", token.Error())
	}

	return nil
}

func (worker *MQTTClient) Subscribe(topic string, handler MqttMessageHandler) error {
	if worker.client == nil {
		return fmt.Errorf("MQTT client is not connected, please start the worker first")
	}

	if _, exists := worker.handlers[topic]; exists {
		fmt.Printf("Handler for topic %s already exists and will be ovewritten\n", topic)
		worker.client.Unsubscribe(topic)
	}

	worker.handlers[topic] = handler

	fmt.Printf("Adding handler for topic %s\n", topic)

	worker.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		handler := worker.handlers[msg.Topic()]
		if handler == nil {
			fmt.Printf("Subscription with no handler registered for topic %s\n", msg.Topic())
			return
		}
		handler(msg)
	})

	return nil
}
