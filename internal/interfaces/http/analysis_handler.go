package http

import (
	"net/http"
	"slices"
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/features/aggregate"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
	"github.com/gin-gonic/gin"
)

type AnalysisHandlerConfig struct {
	AuthorizedDimensions []string `json:"authorized_dimensions"`
}

type AnalysisHandler struct {
	aggregateFeatures   aggregate.AggregateFeatures
	authorizedDimension []string
	log                 *logs.Logger
}

func NewAnalysisHandler(config AnalysisHandlerConfig, aggregateFeatures aggregate.AggregateFeatures, log *logs.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		aggregateFeatures:   aggregateFeatures,
		authorizedDimension: config.AuthorizedDimensions,
		log:                 log,
	}
}

func (h *AnalysisHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/analysis", h.Get)
}

func (h *AnalysisHandler) Get(c *gin.Context) {
	rawDuration, ok := c.GetQuery("duration")
	if !ok {
		c.JSON(http.StatusBadRequest, "Query parameter duration is missing")
		return
	}

	duration, err := time.ParseDuration(rawDuration)
	if err != nil {
		h.log.Error("AnalysisHandler.Get error: can't parse duration", logs.Field{Key: "error", Value: err.Error()})
		c.JSON(http.StatusBadRequest, "Query parameter duration is not in the go time duration format")
		return
	}

	if duration.Seconds() < 0 {
		h.log.Error("AnalysisHandler.Get error: negative duration", logs.Field{Key: "duration", Value: rawDuration})
		c.JSON(http.StatusBadRequest, "Query parameter duration must be a positive value")
		return
	}

	dimension, ok := c.GetQuery("dimension")
	if !ok {
		c.JSON(http.StatusBadRequest, "Query parameter dimension is missing")
		return
	}

	if !slices.Contains(h.authorizedDimension, dimension) {
		h.log.Error("AnalysisHandler.Get error: unauthorized dimension", logs.Field{Key: "supplied_dimension", Value: dimension})
		c.JSON(http.StatusBadRequest, "Unauthorized dimension")
		return
	}

	aggregate, err := h.aggregateFeatures.Aggregate(duration, dimension)
	if err != nil {
		h.log.Error("AnalysisHandler.Get error: ", logs.Field{Key: "error", Value: err.Error()})
		c.JSON(http.StatusInternalServerError, "The server is not able to perform the request")
		return
	}

	c.JSON(http.StatusOK, aggregate)
}
