package sse

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
)

var (
	eventPrefix                     = []byte("data: ")
	ErrReconnectionAttemptsExceeded = errors.New("reconnection attempts exceeded")
)

// Client represents a client for consuming Server-Sent Events (SSE) streams.
// It manages connections to the SSE server, handles reconnections on errors, and broadcasts events to subscribers.
type Client struct {
	url         string
	subscribers []Subscriber

	maxReconnectionAttempts int

	// Mutex to protect Subscribers
	mu sync.Mutex

	closeChan chan struct{}

	log *logs.Logger
}

func NewSSEClient(config Config, log *logs.Logger) *Client {
	return &Client{
		url:                     config.ServerURL,
		subscribers:             []Subscriber{},
		maxReconnectionAttempts: config.MaxReconnectionAttempts,
		mu:                      sync.Mutex{},
		closeChan:               make(chan struct{}),
		log:                     log,
	}
}

// Listen establishes a connection to the SSE server and listens for events in a loop.
// It handles reconnection logic with exponential backoff in case of errors or disconnections.
//
// This function is blocking, it the responsibility of the caller to
// launch it in a go routine.
func (c *Client) Listen() error {
	attempts := 0
	for {
		err := c.readStream()
		if err != nil {
			if attempts > c.maxReconnectionAttempts {
				return ErrReconnectionAttemptsExceeded
			}

			backoff := time.Duration(attempts) * 100 * time.Millisecond

			c.log.Error("SSE Client error, attempting to reconnect to stream",
				logs.Field{Key: "backoff", Value: backoff.String()},
				logs.Field{Key: "error", Value: err.Error()},
			)

			time.Sleep(backoff)
			attempts++
			continue
		}

		// If connection is closed return
		select {
		case <-c.closeChan:
			return nil
		default:
			// Retry connection
		}
	}
}

// NewSubscriber creates and returns a new subscriber.
// It is the caller's responsibility to call RemoveSubscriber to clean up the subscriber
// when it is no longer needed to avoid resource leaks.
func (c *Client) NewSubscriber() (*Subscriber, error) {
	id, err := c.randomID()
	if err != nil {
		return nil, err
	}

	subscriber := Subscriber{
		ID:      id,
		Channel: make(chan []byte),
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.subscribers = append(c.subscribers, subscriber)
	return &subscriber, nil
}

// RemoveSubscriber removes the subscriber with the specified ID from the client's list
// of active subscribers and closes the associated event channel.
// This function must be called by the code that created the subscriber (e.g., after a
// subscriber is no longer needed) to prevent resource leaks.
func (c *Client) RemoveSubscriber(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.subscribers = slices.DeleteFunc(c.subscribers, func(s Subscriber) bool {
		if s.ID == id {
			close(s.Channel)
			return true
		}

		return false
	})
}

func (c *Client) readStream() error {
	res, err := http.Get(c.url)
	if err != nil {
		return fmt.Errorf("can't do request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return &InvalidStatusCodeError{Target: http.StatusOK, Current: res.StatusCode}
	}

	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for {
		select {
		case <-c.closeChan:
			return nil
		default:
			if scanner.Scan() {
				line := scanner.Bytes()
				if !bytes.HasPrefix(line, eventPrefix) {
					continue
				}

				event := bytes.TrimPrefix(line, eventPrefix)

				c.broadcast(event)
				continue
			}

			if scanner.Err() != nil {
				return scanner.Err()
			}
		}
	}
}

func (c *Client) Close() {
	close(c.closeChan)

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, sub := range c.subscribers {
		select {
		case <-sub.Channel:
		default:
			close(sub.Channel)
		}
	}
}

func (c *Client) broadcast(event []byte) {
	// Check if the client is closed
	select {
	case <-c.closeChan:
		return
	default:
	}

	activeSubscribers := make([]Subscriber, 0, len(c.subscribers))

	c.mu.Lock()
	for _, sub := range c.subscribers {
		select {
		case sub.Channel <- event:
			activeSubscribers = append(activeSubscribers, sub)
		default:
		}
	}
	c.subscribers = activeSubscribers
	c.mu.Unlock()
}

func (c *Client) randomID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("can't generate random id")
	}

	return hex.EncodeToString(bytes), nil
}
