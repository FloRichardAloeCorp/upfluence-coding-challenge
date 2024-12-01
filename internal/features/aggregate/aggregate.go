package aggregate

import (
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
)

type AggregateFeatures interface { //nolint:revive
	Aggregate(duration time.Duration, dimension string) (*PostsStatAggregation, error)
}

func NewAggregateFeatures(sseClient *sse.Client) AggregateFeatures {
	repo := &postStatsRepository{
		sseClient: sseClient,
	}

	return newAggregateController(repo)
}
