package config

import (
	"os"
	"slices"
	"testing"
	// model "github.com/FloRichardAloeCorp/upfluence-coding-challenge/pkg/structs"
)

var rawConfig = `{"sse_client_config":{"server_url":"https://stream.upfluence.co/stream","max_reconnection_attempts":10},"router":{"addr":"","port":8080,"gin_mode":"debug","shutdown_timeout":5,"analysis_handler_config":{"authorized_dimensions":["likes","comments","favorites","retweets"]}},"logger":{"level":"INFO"}}`

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	configPath := dir + "/config.json"
	err := os.WriteFile(configPath, []byte(rawConfig), 0666)
	if err != nil {
		t.Fatalf("unexpected error while writing config file: %v", err)
	}

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error while loading config: %v", err)
	}

	if config.Logger.Level != "INFO" {
		t.Errorf("expected Logger.Level to be 'INFO', got '%s'", config.Logger.Level)
	}

	if config.SSEClientConfig.ServerURL != "https://stream.upfluence.co/stream" {
		t.Errorf("expected SSEClientConfig.ServerURL to be 'https://stream.upfluence.co/stream', got '%s'", config.SSEClientConfig.ServerURL)
	}

	if config.SSEClientConfig.MaxReconnectionAttempts != 10 {
		t.Errorf("expected SSEClientConfig.MaxReconnectionAttempts to be 10, got '%d'", config.SSEClientConfig.MaxReconnectionAttempts)
	}

	if config.Router.Addr != "" {
		t.Errorf("expected Router.Addr to be '', got '%s'", config.Router.Addr)
	}

	if config.Router.Port != 8080 {
		t.Errorf("expected Router.Port to be 8080, got '%d'", config.Router.Port)
	}

	if config.Router.GinMode != "debug" {
		t.Errorf("expected Router.GinMode to be 'debug', got '%s'", config.Router.GinMode)
	}

	if config.Router.ShutdownTimeout != 5 {
		t.Errorf("expected Router.ShutdownTimeout to be 5, got '%d'", config.Router.ShutdownTimeout)
	}

	expectedAuthorizedDimensions := []string{
		"likes",
		"comments",
		"favorites",
		"retweets",
	}

	if !slices.Equal(expectedAuthorizedDimensions, config.Router.AnalysisHandlerConfig.AuthorizedDimensions) {
		t.Errorf("invalid Router.AnalysisHandlerConfig.AuthorizedDimensions. expected %v, got %v", expectedAuthorizedDimensions, config.Router.AnalysisHandlerConfig.AuthorizedDimensions)
	}
}

func TestLoadInvalidPath(t *testing.T) {
	_, err := Load("invalid")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestLoadInvalidConfigurationFormat(t *testing.T) {
	dir := t.TempDir()
	configPath := dir + "/config.json"
	err := os.WriteFile(configPath, []byte("this is not a json string"), 0666)
	if err != nil {
		t.Fatalf("unexpected error while writing config file: %v", err)
	}
	_, err = Load(configPath)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
