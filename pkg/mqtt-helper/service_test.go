package mqtthelper

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

	err := Connect(client)
	assert.NoError(t, err)
}

func Test_ServiceConnectTimeout(t *testing.T) {
	client := mocks.NewMockClient(t)
	token := mocks.NewMockToken(t)

	ch := make(chan struct{})

	client.EXPECT().Connect().Return(token).Once()
	token.EXPECT().Done().Return(ch)

	err := Connect(client)
	assert.ErrorIs(t, ErrTimeout, err)
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

	err := Connect(client)
	assert.ErrorIs(t, errFake, err)
}
