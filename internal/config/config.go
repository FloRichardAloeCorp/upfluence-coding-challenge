package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/http"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
)

type Config struct {
	SSEClientConfig sse.Config  `json:"sse_client_config"`
	Router          http.Config `json:"router"`
	Logger          logs.Config `json:"logger"`
}

func Load(path string) (*Config, error) {
	config := new(Config)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't open configuration file: %w", err)
	}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, fmt.Errorf("can't parse configuration: %w", err)
	}

	return config, nil
}
