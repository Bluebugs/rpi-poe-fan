package mqtthelper

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var ErrTimeout = errors.New("timeout")

func Connect(client mqtt.Client) error {
	connection := client.Connect()
	timeout := time.NewTicker(5 * time.Second)
	defer timeout.Stop()

	select {
	case <-connection.Done():
	case <-timeout.C:
		return ErrTimeout
	}

	if connection.Error() != nil {
		return connection.Error()
	}
	return nil
}
