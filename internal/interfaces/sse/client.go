package sse

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"slices"
	"sync"
)

var eventPrefix = []byte("data: ")

type SSEClient struct {
	URL         string
	Subscribers []Subscriber

	// Mutex to protect Subscribers
	mu sync.Mutex

	closeChan chan struct{}
}

func NewSSEClient(url string) *SSEClient {
	return &SSEClient{
		URL:         url,
		Subscribers: []Subscriber{},
		mu:          sync.Mutex{},
		closeChan:   make(chan struct{}),
	}
}

// func (c *SSEClient) Launch() error {
// 	maxAttempts := 5
// 	for i := 0; i < maxAttempts; i++ {
// 		if err := c.Listen(); err != nil {
// 			fmt.Println("Can't listen to sse server, retrying in %d seconds. Error: %w", 1, err)
// 		}

// 		time.Sleep()
// 	}
// }

func (c *SSEClient) Listen() error {
	res, err := http.Get(c.URL)
	if err != nil {
		return fmt.Errorf("can't do request: %w", err)
	}

	go func() {
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		for {
			select {
			case <-c.closeChan:
				return
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
					return
				}
			}
		}
	}()

	return nil
}

func (c *SSEClient) NewSubscriber() (*Subscriber, error) {
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
	c.Subscribers = append(c.Subscribers, subscriber)
	return &subscriber, nil
}

func (c *SSEClient) RemoveSubscriber(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Subscribers = slices.DeleteFunc(c.Subscribers, func(s Subscriber) bool {
		if s.ID == id {
			close(s.Channel)
			return true
		}

		return false
	})
}

func (c *SSEClient) Close() {
	close(c.closeChan)

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, sub := range c.Subscribers {
		select {
		case <-sub.Channel:
		default:
			close(sub.Channel)
		}
	}
}

func (c *SSEClient) broadcast(event []byte) {

	// Check the client is closed
	select {
	case <-c.closeChan:
		return
	default:
	}

	activeSubscribers := make([]Subscriber, 0, len(c.Subscribers))

	c.mu.Lock()
	for _, sub := range c.Subscribers {
		select {
		case sub.Channel <- event:
			activeSubscribers = append(activeSubscribers, sub)
		default:
		}
	}
	c.Subscribers = activeSubscribers
	c.mu.Unlock()
}

func (c *SSEClient) randomID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("can't generate random id")
	}

	return hex.EncodeToString(bytes), nil
}
