//go:build e2e
// +build e2e

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/mocks"
	"github.com/Bluebugs/rpi-poe-fan/pkg/test"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ViableOutput(t *testing.T) {
	page, terminate := newPage(t)
	defer terminate()

	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)
	msg := mocks.NewMockMessage(t)

	s := source{
		client: client,
		rpis:   make(map[string]state),
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
		serve(log, ctx, &s)
		close(shutdown)
	}()

	<-done

	now := time.Now().UTC().Format(time.RFC3339)
	msg.EXPECT().Payload().Return([]byte(fmt.Sprintf(`{"temperature": 50, "fanSpeed": 50, "timestamp": %q}`, now))).Once()
	msg.EXPECT().Topic().Return("/rpi-poe-fan/1/state").Once()
	callback(client, msg)

	// Wait for http service to be running
	time.Sleep(1 * time.Second)

	_, err := page.Goto("http://localhost:8080")
	assert.NoError(t, err)

	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(filepath.Join("testdata", "failed", "screenshot-index.png")),
	})
	assert.NoError(t, err)

	temp := page.Locator("#temp-1")
	err = temp.WaitFor()
	assert.NoError(t, err)

	test.VerifyImage(t, filepath.Join("testdata", "failed", "screenshot-index.png"))

	assert.NoError(t, temp.Err())
	text, err := temp.InnerText()
	assert.NoError(t, err)
	assert.Equal(t, "50.00°C", text)

	fan := page.Locator("#fan-1")
	assert.NoError(t, fan.Err())
	text, err = fan.InnerText()
	assert.NoError(t, err)
	assert.Equal(t, "50%", text)

	cancel()
	<-shutdown
}

func newPage(t *testing.T) (playwright.Page, func()) {
	err := playwright.Install()
	assert.NoError(t, err)

	pw, err := playwright.Run()
	assert.NoError(t, err)
	assert.NotNil(t, pw)

	browser, err := pw.Chromium.Launch()
	assert.NoError(t, err)
	assert.NotNil(t, browser)

	page, err := browser.NewPage()
	assert.NoError(t, err)
	assert.NotNil(t, page)

	err = page.SetViewportSize(1920, 1440)
	assert.NoError(t, err)

	return page, func() {
		page.Close()
		browser.Close()
		pw.Stop()
	}
}