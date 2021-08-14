package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	mqttUrl      string
	scanInterval time.Duration
	ticker       chan bool
}

func NewClient(mqttUrl string, scanInterval time.Duration) *Client {
	return &Client{
		mqttUrl:      mqttUrl,
		scanInterval: scanInterval,
	}
}

func (c *Client) Start() {
	log.Println("Starting client...")

	ticker := c.startTicker(func() {
		updateWebcamSensor()
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

func updateWebcamSensor() {
	webcamActive, err := isWebcamActive()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Webcam active: %t", webcamActive)
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
