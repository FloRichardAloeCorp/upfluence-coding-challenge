package sse

type Config struct {
	ServerURL               string `json:"server_url"`
	MaxReconnectionAttempts int    `json:"max_reconnection_attempts"`
}
