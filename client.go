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

//define a function for the default message handler
var defaultMQTTMessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type Client struct {
	scanInterval time.Duration
	ticker       chan bool
	mqttClient   MQTT.Client
}

func NewClient(mqttBroker string, scanInterval time.Duration) *Client {
	opts := MQTT.NewClientOptions().AddBroker(mqttBroker)
	opts.SetClientID("go-simple") // todo mqtt client id?
	opts.SetDefaultPublishHandler(defaultMQTTMessageHandler)

	//create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &Client{
		scanInterval: scanInterval,
		mqttClient:   mqttClient,
	}
}

func (c *Client) Start() {
	log.Println("Starting client...")

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
		ticker := time.NewTicker(c.scanInterval)
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
	log.Printf("Webcam active: %t", webcamActive)
	c.updateMqttSensor("webcam-active", strconv.FormatBool(webcamActive))
}

func (c *Client) updateMqttSensor(sensorName, value string) {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	topic := fmt.Sprintf("ha-linux-sensors/%s/%s", hostName, sensorName)
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
