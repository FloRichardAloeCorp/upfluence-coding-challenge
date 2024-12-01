package sse

type Subscriber struct {
	ID      string
	Channel chan []byte
}
