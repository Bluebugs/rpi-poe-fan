package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/pkg/event"
	mqtthelper "github.com/Bluebugs/rpi-poe-fan/pkg/mqtt-helper"
	"github.com/gin-gonic/gin"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var entry *template.Template

type state struct {
	Temperature float32
	FanSpeed    uint8
	Timestamp   string
	realtime    time.Time
}

type source struct {
	client mqtt.Client
	rpis   map[string]state
}

func listen(log zerolog.Logger, server string) (*source, error) {
	var err error
	entry, err = template.New("refresh-entry").Funcs(htmlFunc).ParseFS(f, "templates/entry.html", "templates/refresh-entry.html")
	if err != nil {
		return nil, fmt.Errorf("failure to parse entry template: %w", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetClientID(os.Args[0])
	opts.AutoReconnect = true

	client := mqtt.NewClient(opts)
	if err := mqtthelper.Connect(client); err != nil {
		return nil, fmt.Errorf("error establishing connection to mqtt server: %w", err)
	}

	return &source{client: client, rpis: make(map[string]state)}, nil
}

func (s *source) subscribe(sse *event.Event) error {
	token := s.client.Subscribe("/rpi-poe-fan/+/state", 0, func(_ mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		id := topic[len("/rpi-poe-fan/") : len(topic)-len("/state")]

		log.Print("Received message on topic: ", topic)

		var rpi state
		if err := json.Unmarshal(msg.Payload(), &rpi); err != nil {
			log.Error().Err(err).Str("id", id).Msg("Malformed json payload")
			return
		}

		t, err := time.Parse(time.RFC3339, rpi.Timestamp)
		if err != nil {
			log.Error().Err(err).Str("id", id).Str("Timestamp", rpi.Timestamp).Msg("Malformed timestamp in json payload")
			return
		}
		rpi.realtime = t

		if t.Add(5 * time.Minute).Before(time.Now()) {
			log.Error().Str("id", id).Str("Timestamp", rpi.Timestamp).Msg("Timestamp in json payload is too old")
			return
		}

		log.Print("realtime: ", rpi.realtime)

		changed := false
		if current, ok := s.rpis[id]; ok && current.Temperature != rpi.Temperature && current.FanSpeed != rpi.FanSpeed {
			changed = true
		}

		s.rpis[id] = rpi
		log.Info().Str("id", id).Float32("temperature", rpi.Temperature).Uint8("fan", rpi.FanSpeed).Time("timestamp", t).Msg("Registered payload")

		if changed {
			sse.Message <- id
		}
	})
	<-token.Done()

	return token.Error()
}

func (s *source) emit(log zerolog.Logger, id string, c *gin.Context) bool {
	rpi, ok := s.rpis[id]
	if !ok {
		log.Error().Str("Id", id).Msg("Refresh not found")
		return false
	}

	var buf bytes.Buffer
	if err := entry.Execute(&buf, struct {
		Id    string
		State state
	}{id, rpi}); err != nil {
		log.Error().Err(err).Str("Id", id).Msg("Failure to execute entry template")
		return false
	}
	c.SSEvent("Refresh", buf.String())
	return true
}

func (s *source) Close() {
	s.client.Disconnect(0)
}
