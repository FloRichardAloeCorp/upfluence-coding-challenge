package sse

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
)

var loggerInstance, _ = logs.NewLogger(logs.Config{
	Level: "INFO",
})

func createSSEServerMock(interval time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			w.WriteHeader(500)
			return
		}

		for i := 0; i < 5; i++ {
			_, _ = w.Write([]byte("data: dummy event\n\n"))
			flusher.Flush()
			time.Sleep(interval)
		}
	}))
}

func TestSSEClientListen(t *testing.T) {
	server := createSSEServerMock(250 * time.Millisecond)
	defer server.Close()

	client := &Client{
		url: server.URL,
		subscribers: []Subscriber{
			{
				ID:      "id",
				Channel: make(chan []byte),
			},
		},
		closeChan:               make(chan struct{}),
		log:                     loggerInstance,
		maxReconnectionAttempts: 1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var listenError error
	go func() {
		listenError = client.Listen()
	}()

	defer client.Close()

	var receivedEvents [][]byte
	go func() {
		for event := range client.subscribers[0].Channel {
			receivedEvents = append(receivedEvents, event)
		}
	}()

	<-ctx.Done()

	if listenError != nil {
		t.Fatalf("client.Listen returned an error but shouldn't have, error %v", listenError)
	}

	if len(receivedEvents) == 0 {
		t.Fatal("Exepected received events but got 0")
	}

	expectedEventData := []byte("dummy event")
	for _, event := range receivedEvents {
		if !bytes.Equal(event, expectedEventData) {
			t.Fatalf("Expecting event to be %q, got %q", expectedEventData, event)
		}
	}
}

func TestSSEClientListenReconnectionAttempsExceeded(t *testing.T) {
	client := &Client{
		url: "http://dummy.com/stream",
		subscribers: []Subscriber{
			{
				ID:      "id",
				Channel: make(chan []byte),
			},
		},
		closeChan:               make(chan struct{}),
		maxReconnectionAttempts: 2,
		log:                     loggerInstance,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var listenError error
	go func() {
		listenError = client.Listen()
	}()

	defer client.Close()

	var receivedEvents [][]byte
	go func() {
		for event := range client.subscribers[0].Channel {
			receivedEvents = append(receivedEvents, event)
		}
	}()

	<-ctx.Done()

	if listenError == nil {
		t.Fatalf("expected error but have nil")
	}

	if !errors.Is(listenError, ErrReconnectionAttemptsExceeded) {
		t.Fatalf("expected error to be ErrReconnectionAttemptsExceeded got %v", listenError)
	}

	if len(receivedEvents) != 0 {
		t.Fatal("expected 0 received events")
	}
}

func TestSSEClientListenWithoutSubscribers(t *testing.T) {
	server := createSSEServerMock(250 * time.Millisecond)
	defer server.Close()

	client := &Client{
		url:       server.URL,
		closeChan: make(chan struct{}),
		log:       loggerInstance,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer client.Close()

	var listenError error
	go func() {
		listenError = client.Listen()
	}()

	if listenError != nil {
		t.Fatalf("can't listen to sse server: %v", listenError)
	}

	<-ctx.Done()
}

func TestSSEClientClosingClient(t *testing.T) {
	// Creating a mock of the SSE server that sends event during 10 seconds every 2 seconds
	server := createSSEServerMock(2 * time.Second)
	defer server.Close()

	client := &Client{
		url: server.URL,
		subscribers: []Subscriber{
			{
				ID:      "id",
				Channel: make(chan []byte),
			},
		},

		mu:        sync.Mutex{},
		log:       loggerInstance,
		closeChan: make(chan struct{}),
	}

	// Context will end in 1 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var listenError error
	go func() {
		listenError = client.Listen()
	}()

	var receivedEvents [][]byte
	go func() {
		for event := range client.subscribers[0].Channel {
			receivedEvents = append(receivedEvents, event)
		}
	}()

	<-ctx.Done()

	if listenError != nil {
		client.Close()
		t.Fatalf("unexpected error from client.Listen: %v", listenError)
	}

	client.Close()

	if len(receivedEvents) != 1 {
		t.Fatalf("Expected only one received events but got %d", len(receivedEvents))
	}
}

func TestSSEClientNewSubscriber(t *testing.T) {
	client := &Client{
		subscribers: []Subscriber{},
	}

	sub, err := client.NewSubscriber()
	if err != nil {
		t.Fatalf("unexepected error from NewSubscriber: %v", err)
	}

	if sub.ID == "" {
		t.Fatal("Subscriber.ID should not be empty")
	}

	if sub.Channel == nil {
		t.Fatal("Subscriber.Channel should not be nil")
	}

	if len(client.subscribers) != 1 {
		t.Fatalf("Expected client's subscribers len to be %d, got %d", 1, len(client.subscribers))
	}
}

func TestSSEClientRemoveSubscriber(t *testing.T) {
	client := &Client{
		subscribers: []Subscriber{},
	}

	sub, err := client.NewSubscriber()
	if err != nil {
		t.Fatalf("unexepected error from NewSubscriber: %v", err)
	}

	client.RemoveSubscriber(sub.ID)

	if len(client.subscribers) != 0 {
		t.Fatalf("Expected client's subscribers len to be %d, got %d", 0, len(client.subscribers))
	}

	select {
	case _, ok := <-sub.Channel:
		if ok {
			t.Fatalf("expected subscriber channel to be closed")
		}
	default:
		t.Fatalf("channel not closed properly")
	}
}
