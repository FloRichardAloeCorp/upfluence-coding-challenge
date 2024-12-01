package sse

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

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

	client := &SSEClient{
		URL: server.URL,
		Subscribers: []Subscriber{
			{
				ID:      "id",
				Channel: make(chan []byte),
			},
		},
		closeChan: make(chan struct{}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Listen(); err != nil {
		t.Fatalf("can't listen to sse server: %v", err)
	}

	defer client.Close()

	var receivedEvents [][]byte
	go func() {
		for event := range client.Subscribers[0].Channel {
			receivedEvents = append(receivedEvents, event)
		}
	}()

	<-ctx.Done()

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

func TestSSEClientListenWithoutSubscribers(t *testing.T) {
	server := createSSEServerMock(250 * time.Millisecond)
	defer server.Close()

	client := &SSEClient{
		URL:       server.URL,
		closeChan: make(chan struct{}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer client.Close()

	if err := client.Listen(); err != nil {
		t.Fatalf("can't listen to sse server: %v", err)
	}

	<-ctx.Done()
}

func TestSSEClientClosingClient(t *testing.T) {
	// Creating a mock of the SSE server that sends event during 10 seconds every 2 seconds
	server := createSSEServerMock(2 * time.Second)
	defer server.Close()

	client := &SSEClient{
		URL: server.URL,
		Subscribers: []Subscriber{
			{
				ID:      "id",
				Channel: make(chan []byte),
			},
		},

		mu:        sync.Mutex{},
		closeChan: make(chan struct{}),
	}

	// Context will end in 1 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := client.Listen(); err != nil {
		client.Close()
		t.Fatalf("can't listen to sse server: %v", err)
	}

	var receivedEvents [][]byte
	go func() {
		for event := range client.Subscribers[0].Channel {
			receivedEvents = append(receivedEvents, event)
		}
	}()

	<-ctx.Done()
	client.Close()

	if len(receivedEvents) != 1 {
		t.Fatalf("Expected only one received events but got %d", len(receivedEvents))
	}
}

func TestSSEClientNewSubscriber(t *testing.T) {
	client := &SSEClient{
		Subscribers: []Subscriber{},
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

	if len(client.Subscribers) != 1 {
		t.Fatalf("Expected client's subscribers len to be %d, got %d", 1, len(client.Subscribers))
	}
}

func TestSSEClientRemoveSubscriber(t *testing.T) {
	client := &SSEClient{
		Subscribers: []Subscriber{},
	}

	sub, err := client.NewSubscriber()
	if err != nil {
		t.Fatalf("unexepected error from NewSubscriber: %v", err)
	}

	client.RemoveSubscriber(sub.ID)

	if len(client.Subscribers) != 0 {
		t.Fatalf("Expected client's subscribers len to be %d, got %d", 0, len(client.Subscribers))
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
