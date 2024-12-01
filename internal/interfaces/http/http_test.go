package http

import (
	"testing"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
)

func TestNewRouter(t *testing.T) {
	log, err := logs.NewLogger(logs.Config{
		Level: "INFO",
	})
	if err != nil {
		t.Fatal("unexpected error during logger instanciation")
	}

	conf := Config{
		GinMode:         "debug",
		Addr:            "",
		Port:            8080,
		ShutdownTimeout: 10,
	}

	router := NewRouter(conf, log)
	if router == nil {
		t.Error("instanciated router should not be nil")
	}
}
