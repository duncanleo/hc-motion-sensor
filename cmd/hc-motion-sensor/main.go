package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/gjson"
)

var (
	motion = false
)

func connect(clientID string, uri *url.URL) (mqtt.Client, error) {
	var opts = mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientID)
	opts.SetCleanSession(false)

	var client = mqtt.NewClient(opts)
	var token = client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	return client, token.Error()
}

func main() {
	var name = flag.String("name", "hc-motion-sensor", "name of the sensor to display in HomeKit")
	var manufacturer = flag.String("manufacturer", "TZT", "manufacturer of the sensor")
	var model = flag.String("model", "FC-37", "model number of the sensor")
	var serial = flag.String("serial", "0000", "serial number of the sensor")
	var pin = flag.String("pin", "00102003", "PIN number to pair the HomeKit accessory")
	var port = flag.String("port", "", "Port number for the HomeKit accessory")
	var storagePath = flag.String("storagePath", "hc-motion-sensor-data", "path to store data")

	var brokerURI = flag.String("brokerURI", "mqtt://127.0.0.1:1883", "URI of the MQTT broker")
	var clientID = flag.String("clientID", "hc-motion-sensor", "client ID for MQTT")

	var topic = flag.String("topic", "motion", "topic to subscribe in MQTT")
	var jsonPath = flag.String("jsonPath", "wet", "path to JSON boolean")

	flag.Parse()

	mqttURI, err := url.Parse(*brokerURI)
	if err != nil {
		log.Fatal(err)
	}

	info := accessory.Info{
		Name:         *name,
		Manufacturer: *manufacturer,
		Model:        *model,
		SerialNumber: *serial,
	}

	ac := accessory.New(info, accessory.TypeSensor)

	motionSensor := service.NewMotionSensor()
	ac.AddService(motionSensor.Service)

	motionSensor.MotionDetected.SetValue(false)
	motionSensor.MotionDetected.OnValueGet(func() interface{} { return motion })

	client, err := connect(*clientID, mqttURI)
	if err != nil {
		log.Fatal(err)
	}

	client.Subscribe(*topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		motion = gjson.Get(string(msg.Payload()), *jsonPath).Bool()
		motionSensor.MotionDetected.UpdateValue(motion)
	})

	hcConfig := hc.Config{
		Pin:         *pin,
		StoragePath: *storagePath,
		Port:        *port,
	}

	t, err := hc.NewIPTransport(hcConfig, ac)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

}
