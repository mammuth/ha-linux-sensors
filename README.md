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

### Run on boot with systemd

Create a systemd service at `/etc/systemd/system/ha-linux-sensors.service`:

> Replace `/home/max/Code/ha-linux-sensors` in the example below with the absolute path to your local checkout.

```ini
[Unit]
Description=HA Linux Sensors
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=max
WorkingDirectory=/home/max/Code/ha-linux-sensors
EnvironmentFile=/etc/ha-linux-sensors.env
ExecStart=/home/max/Code/ha-linux-sensors/ha-linux-sensors -mqttBroker ${MQTT_BROKER} -mqttUser ${MQTT_USER} -mqttPassword ${MQTT_PASSWORD}
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Create the env file at `/etc/ha-linux-sensors.env`:

```bash
MQTT_BROKER=tcp://192.168.178.20:1883
MQTT_USER=your-mqtt-username
MQTT_PASSWORD=replace-with-your-password
```

Set secure permissions and enable the service:

```bash
sudo chmod 600 /etc/ha-linux-sensors.env
sudo systemctl daemon-reload
sudo systemctl enable --now ha-linux-sensors.service
sudo systemctl status ha-linux-sensors.service
```

At the Home Assistant side, you probably want to create a [MQTT sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt/).
