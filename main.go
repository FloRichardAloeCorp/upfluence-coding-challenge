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
	ConfigPathEnvVar  = "API_CONFIG"
	DefaultConfigPath = "./config.json"
)

func main() {
	configFilePath, present := os.LookupEnv(ConfigPathEnvVar)
	if !present {
		configFilePath = DefaultConfigPath
	}

	config, err := config.Load(configFilePath)
	if err != nil {
		panic(err)
	}

	log, err := logs.NewLogger(config.Logger)
	if err != nil {
		panic(err)
	}

	run, shutdown, err := app.Launch(*config, log)
	if err != nil {
		panic(err)
	}

	go run()

	WaitSignalShutdown(shutdown, log)
}

func WaitSignalShutdown(shutdown func() error, log *logs.Logger) {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	if err := shutdown(); err != nil {
		log.Error("Fail to properly close the server", logs.Field{Key: "error", Value: err.Error()})
	}
}
