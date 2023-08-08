package main

import (
	"errors"
	"testing"

	"github.com/Bluebugs/rpi-poe-fan/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_ServiceConnect(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)

	ch := make(chan struct{})
	close(ch)

	client.EXPECT().Connect().Return(token).Once()
	token.EXPECT().Done().Return(ch)
	token.EXPECT().Error().Return(nil)

	err := connect(client)
	assert.NoError(t, err)
}

func Test_ServiceConnectTimeout(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)

	ch := make(chan struct{})

	client.EXPECT().Connect().Return(token).Once()
	token.EXPECT().Done().Return(ch)

	err := connect(client)
	assert.ErrorIs(t, errTimeout, err)
}

func Test_ServiceConnectError(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)

	ch := make(chan struct{})
	close(ch)

	errFake := errors.New("fake error")

	client.EXPECT().Connect().Return(token).Once()
	token.EXPECT().Done().Return(ch)
	token.EXPECT().Error().Return(errFake)

	err := connect(client)
	assert.ErrorIs(t, errFake, err)
}
