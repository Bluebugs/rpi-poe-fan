package main

import (
	"context"
	"errors"
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
	f, err := fans.HwMon()
	if err != nil {
		log.Fatal("Failure to list PoE Hat fan:", err)
	}

	t := cpu.NewRPiTemp()

	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		interval = 3 * time.Second
	}

	id, err := machineid.ProtectedID(os.Args[0])
	if err != nil {
		log.Fatal("Failure to get machine id:", err)
	}

	if err := serve("tcp://localhost:1883", id, f, t,
		func() { _, _ = daemon.SdNotify(false, daemon.SdNotifyReady) },
		func() error {
			time.Sleep(interval / 3)
			_, _ = daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			return nil
		}); err != nil {
		log.Fatal("Failure to serve:", err)
	}
}

func serve(server string, id string, f fans.Fan, t cpu.Temp, ready func(), tick func() error) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetClientID("sensor" + os.Args[0])
	opts.AutoReconnect = true

	client := mqtt.NewClient(opts)
	if err := mqtthelper.Connect(client); err != nil {
		return err
	}
	defer client.Disconnect(0)

	if err := subscribe(client, id, f); err != nil {
		return err
	}

	ready()
	for {
		if err := publish(client, id, t, f); err != nil {
			log.Println(err)
		}

		if err := tick(); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
}
