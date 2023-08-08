package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Bluebugs/rpi-poe-fan/mocks"
	"github.com/Bluebugs/rpi-poe-fan/pkg/cpu"
	"github.com/Bluebugs/rpi-poe-fan/pkg/fans"
	"github.com/stretchr/testify/assert"
)

func Test_Publish(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)
	temp := cpu.NewMockTemp(t)
	fan := fans.NewMockFan(t)

	ch := make(chan struct{})
	close(ch)

	now = func() time.Time {
		return time.Date(2023, 8, 7, 6, 5, 4, 0, time.UTC)
	}

	token.EXPECT().Done().Return(ch)
	token.EXPECT().Error().Return(nil)

	errFan := errors.New("fan error")
	errCPU := errors.New("cpu error")

	tests := map[string]struct {
		cpuTemp       float32
		cpuError      error
		fanSpeed      uint8
		fanError      error
		expectedError error
	}{
		"hot":       {cpuTemp: 70, fanSpeed: 100},
		"cold":      {cpuTemp: 30, fanSpeed: 0},
		"medium":    {cpuTemp: 50, fanSpeed: 50},
		"fan error": {cpuTemp: 42, fanSpeed: 7, fanError: errFan, expectedError: errFan},
		"cpu error": {cpuTemp: 35, fanSpeed: 24, cpuError: errCPU, expectedError: errCPU},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			temp.EXPECT().Read().Return(tc.cpuTemp, tc.cpuError).Once()
			if tc.cpuError == nil {
				fan.EXPECT().Speed().Return(tc.fanSpeed, tc.fanError).Once()
			}

			if tc.expectedError == nil {
				expectedJSON := fmt.Sprintf("{\"Temperature\":%d,\"FanSpeed\":%d,\"Timestamp\":\"2023-08-07T06:05:04Z\"}", int(tc.cpuTemp), int(tc.fanSpeed))
				client.EXPECT().Publish("/rpi-poe-fan/test/state", byte(0), true, expectedJSON).Return(token).Once()
			}

			err := publish(client, "test", temp, fan)
			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
