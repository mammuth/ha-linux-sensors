package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var defaultMQTTMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type ClientConfig struct {
	scanInterval time.Duration
	mqttBroker   string
	mqttUser     string
	mqttPassword string
}
type Client struct {
	config     *ClientConfig
	ticker     chan bool
	mqttClient MQTT.Client
}

func NewClient(config *ClientConfig) *Client {
	log.Printf("Using MQTT broker %s", config.mqttBroker)
	log.Printf("Scan interval: %s seconds", strconv.Itoa(int(config.scanInterval.Seconds())))

	opts := MQTT.NewClientOptions().AddBroker(config.mqttBroker)
	opts.SetClientID("ha-linux-sensors")
	opts.SetUsername(config.mqttUser)
	opts.SetPassword(config.mqttPassword)
	opts.SetDefaultPublishHandler(defaultMQTTMessageHandler)

	// todo what happens when the connection is lsot? how to reconenct?
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		msg := fmt.Sprintf("Error connecting to broker %s", mqttBroker)
		fmt.Errorf(msg, token.Error())
	}

	return &Client{
		config:     config,
		mqttClient: mqttClient,
	}
}

func (c *Client) Start() {
	log.Println("Starting routine...")

	ticker := c.startTicker(func() {
		c.updateWebcamSensor()
	})
	c.ticker = ticker
}

func (c *Client) Stop() {
	close(c.ticker)
}

func (c *Client) startTicker(f func()) chan bool {
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(c.config.scanInterval)
		defer ticker.Stop()
		for {
			f()
			select {
			case <-ticker.C:
				continue
			case <-done:
				fmt.Println("done")
				return
			}
		}
	}()
	return done
}

func (c *Client) updateWebcamSensor() {
	webcamActive, err := isWebcamActive()
	if err != nil {
		log.Println(err)
		return
	}
	sensorValue := "off"
	if webcamActive {
		sensorValue = "on"
	}
	c.updateMqttSensor("webcam", sensorValue)
}

func (c *Client) updateMqttSensor(sensorName, value string) {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	topic := fmt.Sprintf("ha-linux-sensors/%s/%s", hostName, sensorName)
	log.Printf("Publishing sensor update %s:%s", topic, value)
	token := c.mqttClient.Publish(topic, 0, false, value)
	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			log.Printf("Error publishing MQTT topic: %s", token.Error())
		}
	}()
}

func isWebcamActive() (bool, error) {
	cmd := "lsmod | grep uvcvideo | awk '{print $NF}' | head -n1"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return false, fmt.Errorf("Failed getting webcam status", cmd, err)
	}
	stringResult := strings.TrimSuffix(string(out), "\n")
	result, err := strconv.Atoi(stringResult)
	if err != nil {
		return false, fmt.Errorf("Failed getting webcam status", err)
	}
	if result != 0 && result != 1 {
		return false, errors.New("webcam bash command return an unexpected value")
	}
	return result == 1, nil
}
