package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/config"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/features/aggregate"
	ginhttp "github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/http"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
)

type RunCallback func()
type CloseCallback func() error

func Launch(config config.Config, log *logs.Logger) (RunCallback, CloseCallback, error) {
	sseClient := sse.NewSSEClient(config.SSEClientConfig.ServerURL)

	aggregateFeature := aggregate.NewAggregateFeatures(sseClient)

	router := ginhttp.NewRouter(config.Router, log)

	analysisHandler := ginhttp.NewAnalysisHandler(config.Router.AnalysisHandlerConfig, aggregateFeature, log)

	analysisHandler.RegisterRoutes(router)

	addrGin := config.Router.Addr + ":" + strconv.Itoa(config.Router.Port)
	srv := &http.Server{
		ReadHeaderTimeout: time.Millisecond,
		Addr:              addrGin,
		Handler:           router,
	}

	close := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Router.ShutdownTimeout)*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("can't shutdown server: %w", err)
		}

		sseClient.Close()

		return nil
	}

	run := func() {
		sseClient.Listen()
		log.Info("REST API listening on " + addrGin)
		log.Error(router.Run(addrGin).Error())
	}

	return run, close, nil
}
