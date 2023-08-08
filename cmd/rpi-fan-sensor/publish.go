package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/pkg/cpu"
	"github.com/Bluebugs/rpi-poe-fan/pkg/fans"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type state struct {
	Temperature float32
	FanSpeed    uint8
	Timestamp   string
}

var now = time.Now

func publish(client mqtt.Client, id string, t cpu.Temp, fan fans.Fan) error {
	temp, err := t.Read()
	if err != nil {
		return err
	}

	speed, err := fan.Speed()
	if err != nil {
		return err
	}

	s := state{
		Temperature: temp,
		FanSpeed:    speed,
		Timestamp:   now().Format(time.RFC3339),
	}

	js, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Publishing CPU temperature:", temp, "°C", "Fan speed at", speed, "%")

	token := client.Publish("/rpi-poe-fan/"+id+"/state", 0, true, string(js))
	<-token.Done()

	return token.Error()
}
