package main

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var errTimeout = errors.New("timeout")

func connect(client mqtt.Client) error {
	connection := client.Connect()
	timeout := time.NewTicker(5 * time.Second)
	defer timeout.Stop()

	select {
	case <-connection.Done():
	case <-timeout.C:
		return errTimeout
	}

	if connection.Error() != nil {
		return connection.Error()
	}
	return nil
}
