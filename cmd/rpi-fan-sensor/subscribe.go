package main

import (
	"log"
	"strconv"

	"github.com/Bluebugs/rpi-poe-fan/pkg/fans"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func subscribe(client mqtt.Client, id string, fan fans.Fan) error {
	topic := "/rpi-poe-fan/" + id + "/speed"
	token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		speed, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			log.Println(err)
			return
		}

		if speed > 100 || speed < 0 {
			log.Println("Speed must be a percentage between 0 and 100.")
			return
		}

		if err := fan.SetSpeed(uint8(speed)); err != nil {
			log.Println(err)
			return
		}

		log.Println("Fan speed set to", speed, "%")
	})
	<-token.Done()

	return token.Error()
}
