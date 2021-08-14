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
	mqttBroker   = flag.String("mqttBroker", "", "URI of the MQTT broker, eg. tcp://mqtt.eclipseprojects.io:1883")
	scanInterval = flag.Int("interval", 10, "Scan interval in seconds")
)

func main() {

	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	client := NewClient(*mqttBroker, time.Duration(*scanInterval)*time.Second)

	client.Start()

	<-sigs
	log.Println("Received interrupt, shutting down")
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client.Stop()
}
