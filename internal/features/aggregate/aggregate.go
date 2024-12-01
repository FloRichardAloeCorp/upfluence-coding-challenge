package aggregate

import (
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
)

type AggregateFeatures interface {
	Aggregate(duration time.Duration, dimension string) (*PostsStatAggregation, error)
}

func NewAggregateFeatures(sseClient *sse.SSEClient) AggregateFeatures {
	repo := &postStatsRepository{
		sseClient: sseClient,
	}

	return newAggregateController(repo)
}
