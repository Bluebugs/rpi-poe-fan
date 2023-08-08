package main

import (
	"testing"

	"github.com/Bluebugs/rpi-poe-fan/mocks"
	"github.com/Bluebugs/rpi-poe-fan/pkg/fans"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Subscribe(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)
	msg := mocks.NewMockMessage(t)
	fan := fans.NewMockFan(t)

	var callback mqtt.MessageHandler

	client.EXPECT().Subscribe("/rpi-poe-fan/test/speed", uint8(0), mock.Anything).RunAndReturn(func(s string, b byte, mh mqtt.MessageHandler) mqtt.Token {
		assert.Equal(t, "/rpi-poe-fan/test/speed", s)
		assert.Equal(t, byte(0), b)
		assert.NotNil(t, mh)

		callback = mh

		ch := make(chan struct{})
		defer close(ch)

		token.EXPECT().Done().Return(ch)
		token.EXPECT().Error().Return(nil)
		return token
	})

	if err := subscribe(client, "test", fan); err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, callback)

	msg.EXPECT().Payload().Return([]byte("50")).Once()
	fan.EXPECT().SetSpeed(uint8(50)).Return(nil).Once()
	callback(client, msg)

	msg.EXPECT().Payload().Return([]byte("0")).Once()
	fan.EXPECT().SetSpeed(uint8(0)).Return(nil).Once()
	callback(client, msg)

	msg.EXPECT().Payload().Return([]byte("100")).Once()
	fan.EXPECT().SetSpeed(uint8(100)).Return(nil).Once()
	callback(client, msg)
}

func Test_SubscribeMsgError(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)
	msg := mocks.NewMockMessage(t)
	fan := fans.NewMockFan(t)

	var callback mqtt.MessageHandler

	client.EXPECT().Subscribe("/rpi-poe-fan/test/speed", uint8(0), mock.Anything).RunAndReturn(func(s string, b byte, mh mqtt.MessageHandler) mqtt.Token {
		assert.Equal(t, "/rpi-poe-fan/test/speed", s)
		assert.Equal(t, byte(0), b)
		assert.NotNil(t, mh)

		callback = mh

		ch := make(chan struct{})
		defer close(ch)

		token.EXPECT().Done().Return(ch)
		token.EXPECT().Error().Return(nil)
		return token
	})

	if err := subscribe(client, "test", fan); err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, callback)

	msg.EXPECT().Payload().Return([]byte("nan")).Once()
	callback(client, msg)

	msg.EXPECT().Payload().Return([]byte("-1")).Once()
	callback(client, msg)

	msg.EXPECT().Payload().Return([]byte("101")).Once()
	callback(client, msg)
}
