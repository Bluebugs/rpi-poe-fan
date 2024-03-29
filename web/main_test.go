package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/mocks"
	"github.com/Bluebugs/rpi-poe-fan/web/types"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/graceful"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_JSONEndpoint(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)
	msg := mocks.NewMockMessage(t)

	s := source{
		client: client,
		rpis:   make(map[string]types.State),
	}

	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stderr
		w.TimeFormat = time.RFC3339
	})
	log := zerolog.New(writer).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &log

	ctx, cancel := context.WithCancel(context.Background())

	var callback mqtt.MessageHandler
	done := make(chan struct{})
	shutdown := make(chan struct{})

	client.EXPECT().Subscribe("/rpi-poe-fan/+/state", uint8(0), mock.Anything).RunAndReturn(func(s string, b byte, mh mqtt.MessageHandler) mqtt.Token {
		assert.Equal(t, "/rpi-poe-fan/+/state", s)
		assert.Equal(t, byte(0), b)
		assert.NotNil(t, mh)

		callback = mh
		defer close(done)

		ch := make(chan struct{})
		defer close(ch)

		token.EXPECT().Done().Return(ch)
		token.EXPECT().Error().Return(nil)
		return token
	})
	client.EXPECT().Disconnect(uint(0)).Return().Once()

	go func() {
		_ = Serve(&log, ctx, &s, func() error { return nil }, graceful.WithAddr("localhost:9080"))
		close(shutdown)
	}()

	<-done

	now := time.Now().UTC().Format(time.RFC3339)
	realtime, _ := time.Parse(time.RFC3339, now)
	msg.EXPECT().Payload().Return([]byte(fmt.Sprintf(`{"temperature": 50, "fanSpeed": 50, "timestamp": %q}`, now))).Once()
	msg.EXPECT().Topic().Return("/rpi-poe-fan/1/state").Once()
	callback(client, msg)

	time.Sleep(10 * time.Millisecond)

	rpis := map[string]types.State{}

	err := httpGet("http://localhost:9080/api/entries", &rpis)
	assert.NoError(t, err)

	assert.Len(t, rpis, 1)
	assert.Equal(t, types.State{Temperature: 50, FanSpeed: 50, Timestamp: now, Realtime: realtime}, rpis["1"])

	st := types.State{}
	err = httpGet("http://localhost:9080/api/entry/1", &st)
	assert.NoError(t, err)

	assert.Equal(t, types.State{Temperature: 50, FanSpeed: 50, Timestamp: now, Realtime: realtime}, st)

	cancel()
	<-shutdown
}

func httpGet(target string, object interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-type") != "application/json; charset=utf-8" {
		return fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, object); err != nil {
		return err
	}

	return nil
}
