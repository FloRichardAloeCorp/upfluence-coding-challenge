package http

import (
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/http/middlewares"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
	"github.com/gin-gonic/gin"
)

type Config struct {
	GinMode               string                `json:"gin_mode"`
	Addr                  string                `json:"addr"`
	Port                  int                   `json:"port"`
	ShutdownTimeout       int                   `json:"shutdown_timeout"`
	AnalysisHandlerConfig AnalysisHandlerConfig `json:"analysis_handler_config"`
}

func NewRouter(config Config, log *logs.Logger) *gin.Engine {
	router := gin.New()
	gin.SetMode(config.GinMode)

	router.Use(gin.Recovery())
	router.Use(middlewares.RequestsLogger(log))

	return router
}
