package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/app"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/config"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
)

const (
	PREFIX_ENV          = "API"
	ENV_CONFIG          = PREFIX_ENV + "_CONFIG"
	DEFAULT_CONFIG_PATH = "./config.json"
)

func main() {
	configFilePath, present := os.LookupEnv(ENV_CONFIG)
	if !present {
		configFilePath = DEFAULT_CONFIG_PATH
	}

	config, err := config.Load(configFilePath)
	if err != nil {
		panic(err)
	}

	log, err := logs.NewLogger(config.Logger)
	if err != nil {
		panic(err)
	}

	run, close, err := app.Launch(*config, log)
	if err != nil {
		panic(err)
	}

	go run()

	WaitSignalShutdown(close, log)
}

func WaitSignalShutdown(close func() error, log *logs.Logger) {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	close()
}
