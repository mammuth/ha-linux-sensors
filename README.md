# ha-linux-sensors

A go program that sends sensor data from Linux machines to Home Assistant via MQTT (WIP)

Currently supported sensors:
- Webcam enabled (eg. `update ha-linux-sensors/<hostname>/webcam:on`)

![image](https://user-images.githubusercontent.com/3121306/129533094-9eb04987-541f-4809-a9fa-0068e0bd96b1.png)


### Usage

Clone and build it with `go build` (or download a potentially outdated binary from [here](https://github.com/mammuth/ha-linux-sensors/releases))

```
$ ./ha-linux-sensors -h
Usage of ./ha-linux-sensors:
  -interval int
    	Scan interval in seconds (default 10)
  -mqttBroker string
    	URI of the MQTT broker, eg. tcp://broker.hivemq.com:1883
  -mqttPassword string
    	Password for the mqtt connection
  -mqttUser string
    	Username for the mqtt connection
```

You can then add the command to your startup tools.

At the Home Assistant side, you probably want to create a [MQTT sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt/).
