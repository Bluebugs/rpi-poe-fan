package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/pkg/cpu"
	"github.com/Bluebugs/rpi-poe-fan/pkg/fans"
	"github.com/Bluebugs/rpi-poe-fan/pkg/test"
	"github.com/Bluebugs/rpi-poe-fan/web"

	"github.com/gin-contrib/graceful"
	"github.com/neilotoole/errgroup"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const end2endServerAddr = "localhost:9981"

func Test_End2End(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	page, terminate := newPage(t)
	defer terminate()

	mqttServer, cleanup, err := setupMqttServer(ctx)
	if err != nil {
		t.Fatalf("Failed to setup MQTT server: %v", err)
	}
	defer cleanup()

	require.NotEmpty(t, mqttServer)

	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = os.Stderr
		w.TimeFormat = time.RFC3339
	})
	log := zerolog.New(writer).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &log

	source, err := web.Listen(mqttServer)
	require.NoError(t, err)
	require.NotNil(t, source)

	temp := cpu.NewMockTemp(t)
	fan := fans.NewMockFan(t)

	stop := make(chan struct{})

	var eg errgroup.Group
	eg.Go(func() error {
		return serve(mqttServer, "1", fan, temp, func() { <-stop }, func() { <-stop })
	})

	eg.Go(func() error {
		return web.Serve(&log, ctx, source, func() error { return nil }, graceful.WithAddr(end2endServerAddr))
	})

	// Wait for http service to be running
	time.Sleep(1 * time.Second)

	tests := map[string]struct {
		cpuTemp  float32
		fanSpeed uint8
	}{
		"hot":    {cpuTemp: 70, fanSpeed: 100},
		"cold":   {cpuTemp: 30, fanSpeed: 0},
		"medium": {cpuTemp: 50, fanSpeed: 50},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			temp.EXPECT().Read().Return(tc.cpuTemp, nil).Once()
			fan.EXPECT().Speed().Return(tc.fanSpeed, nil).Once()

			stop <- struct{}{}

			// Wait for message to propagate
			time.Sleep(16 * time.Millisecond)

			_, err := page.Goto("http://" + end2endServerAddr)
			assert.NoError(t, err)

			temp := page.Locator("#temp-1")
			err = temp.WaitFor()
			assert.NoError(t, err)
			assert.NoError(t, temp.Err())
			text, err := temp.InnerText()
			assert.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("%.02f°C", tc.cpuTemp), text)

			fan := page.Locator("#fan-1")
			assert.NoError(t, fan.Err())
			text, err = fan.InnerText()
			assert.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("%d%%", tc.fanSpeed), text)

			_, err = page.Screenshot(playwright.PageScreenshotOptions{
				Path: playwright.String(filepath.Join("testdata", "failed", "screenshot-"+name+".png")),
			})
			assert.NoError(t, err)
			test.VerifyImage(t, filepath.Join("testdata", "failed", "screenshot-"+name+".png"))
		})
	}
}

func setupMqttServer(ctx context.Context) (string, func(), error) {
	// Create a new container
	req := testcontainers.ContainerRequest{
		Image:        "eclipse-mosquitto",
		ExposedPorts: []string{"1883/tcp"},
		WaitingFor:   wait.ForListeningPort("1883/tcp"),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      filepath.Join("testdata", "mosquitto.conf"),
				ContainerFilePath: "/mosquitto/config/mosquitto.conf",
				FileMode:          0o600,
			},
		},
	}

	mqttContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", func() {}, fmt.Errorf("Failed to create MQTT container: %v", err)
	}
	cancel := true
	defer func() {
		if cancel {
			mqttContainer.Terminate(ctx)
		}
	}()

	mqttContainerIp, err := mqttContainer.Host(ctx)
	if err != nil {
		return "", func() {}, fmt.Errorf("Failed to get MQTT container IP: %v", err)
	}

	mqttContainerPort, err := mqttContainer.MappedPort(ctx, "1883")
	if err != nil {
		return "", func() {}, fmt.Errorf("Failed to get MQTT container port: %v", err)
	}

	cancel = false
	return fmt.Sprintf("tcp://%s:%s", mqttContainerIp, mqttContainerPort.Port()), func() {
		stream, err := mqttContainer.CopyFileFromContainer(ctx, "/mosquitto/log/mosquitto.log")
		if err == nil {
			defer stream.Close()
			b, _ := io.ReadAll(stream)
			os.WriteFile("mosquitto.log", b, 0o600)
		}
		mqttContainer.Terminate(ctx)
	}, nil
}

func newPage(t *testing.T) (playwright.Page, func()) {
	err := playwright.Install()
	assert.NoError(t, err)

	pw, err := playwright.Run()
	assert.NoError(t, err)
	assert.NotNil(t, pw)

	browser, err := pw.Firefox.Launch()
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
