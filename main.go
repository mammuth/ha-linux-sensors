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
	mqttUrl      = flag.String("mqttUrl", "https://localhost", "URL of the MQTT server")
	scanInterval = flag.Int("interval", 10, "Scan interval in seconds")
)

func main() {

	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	client := NewClient(*mqttUrl, time.Duration(*scanInterval)*time.Second)

	client.Start()

	<-sigs
	log.Println("Received interrupt, shutting down")
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client.Stop()
}
