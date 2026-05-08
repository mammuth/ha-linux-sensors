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

### Run in background with systemd

Create a user-level systemd service at `~/.config/systemd/user/ha-linux-sensors.service`:

> Replace `%h/Code/ha-linux-sensors` in the example below with the absolute path to your local checkout if needed.

```ini
[Unit]
Description=HA Linux Sensors
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=%h/Code/ha-linux-sensors
EnvironmentFile=%h/.config/ha-linux-sensors/env
ExecStart=%h/Code/ha-linux-sensors/ha-linux-sensors -mqttBroker ${MQTT_BROKER} -mqttUser ${MQTT_USER} -mqttPassword ${MQTT_PASSWORD}
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
```

Create the env file at `~/.config/ha-linux-sensors/env`:

```bash
MQTT_BROKER=tcp://192.168.178.20:1883
MQTT_USER=your-mqtt-username
MQTT_PASSWORD=replace-with-your-password
```

Set secure permissions and enable the service:

```bash
mkdir -p ~/.config/systemd/user ~/.config/ha-linux-sensors
chmod 600 ~/.config/ha-linux-sensors/env
systemctl --user daemon-reload
systemctl --user enable --now ha-linux-sensors.service
systemctl --user status ha-linux-sensors.service
```

Optional: run at boot before login (requires one-time root command):

```bash
sudo loginctl enable-linger <your-user>
```

### HomeAssistant integration
You probably want to create a [MQTT sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt/).
