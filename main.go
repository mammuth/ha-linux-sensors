package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	scanInterval = flag.Int("interval", 10, "Scan interval in seconds")
	mqttBroker   = flag.String("mqttBroker", "", "URI of the MQTT broker, eg. tcp://broker.hivemq.com:1883")
	mqttUser     = flag.String("mqttUser", "", "Username for the mqtt connection")
	mqttPassword = flag.String("mqttPassword", "", "Password for the mqtt connection")
)

func main() {

	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	client := NewClient(&ClientConfig{
		scanInterval: time.Duration(*scanInterval) * time.Second,
		mqttBroker:   *mqttBroker,
		mqttUser:     *mqttUser,
		mqttPassword: *mqttPassword,
	})

	client.Start()

	<-sigs
	log.Println("Received interrupt, shutting down")
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client.Stop()
}
