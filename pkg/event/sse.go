package event

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan string

// Initialize event and Start procnteessing requests
func New() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	log := zerolog.DefaultContextLogger

	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Info().Int("TotalClients", len(stream.TotalClients)).Msg("Client added.")

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Info().Int("TotalClients", len(stream.TotalClients)).Msg("Removed client.")

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *Event) GetChannel(c *gin.Context) (ClientChan, error) {
	v, ok := c.Get("clientChan")
	if !ok {
		return nil, errors.New("clientChan not found")
	}
	clientChan, ok := v.(ClientChan)
	if !ok {
		return nil, errors.New("clientChan is not ClientChan")
	}
	return clientChan, nil
}

func (stream *Event) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
