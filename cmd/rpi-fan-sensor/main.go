package main

import (
	"log"
	"os"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/pkg/cpu"
	"github.com/Bluebugs/rpi-poe-fan/pkg/fans"
	mqtthelper "github.com/Bluebugs/rpi-poe-fan/pkg/mqtt-helper"
	"github.com/coreos/go-systemd/daemon"

	"github.com/denisbrodbeck/machineid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	id, err := machineid.ProtectedID(os.Args[0])
	if err != nil {
		log.Fatal("Failure to find machine unique ID:", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID(os.Args[0])
	opts.AutoReconnect = true

	client := mqtt.NewClient(opts)
	if err := mqtthelper.Connect(client); err != nil {
		log.Fatal("Failure to establish local MQTT connection:", err)
	}
	defer client.Disconnect(0)

	f, err := fans.HwMon()
	if err != nil {
		log.Fatal("Failure to list PoE Hat fan:", err)
	}

	t := cpu.NewRPiTemp()

	if err := subscribe(client, id, f); err != nil {
		log.Fatal("Failure to subscribe to MQTT topics:", err)
	}

	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		interval = 3 * time.Second
	}

	daemon.SdNotify(false, daemon.SdNotifyReady)
	for {
		if err := publish(client, id, t, f); err != nil {
			log.Println(err)
		}

		time.Sleep(interval / 3)
		daemon.SdNotify(false, daemon.SdNotifyWatchdog)
	}
}
