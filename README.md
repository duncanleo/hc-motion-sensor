# hc-rain-sensor
This is a HomeKit rain sensor built with [hc](https://github.com/brutella/hc). It emulates a motion sensor.

### Usage
```shell
Usage of hc-rain-sensor:
  -brokerURI string
        URI of the MQTT broker (default "mqtt://127.0.0.1:1883")
  -clientID string
        client ID for MQTT (default "hc-rain-sensor")
  -jsonPath string
        path to JSON boolean (default "wet")
  -manufacturer string
        manufacturer of the sensor (default "TZT")
  -model string
        model number of the sensor (default "FC-37")
  -name string
        name of the sensor to display in HomeKit (default "hc-rain-sensor")
  -pin string
        PIN number to pair the HomeKit accessory (default "00102003")
  -port string
        Port number for the HomeKit accessory
  -serial string
        serial number of the sensor (default "0000")
  -storagePath string
        path to store data (default "hc-rain-sensor-data")
  -topic string
        topic to subscribe in MQTT (default "rain")
```

### JSON Path
The code uses the [gjson](https://github.com/tidwall/gjson) package to parse data freely from any JSON response. The key system is similar to `jq` but it omits the leading period (`.`). See this [playground](http://tidwall.com/gjson-play) for more info.

