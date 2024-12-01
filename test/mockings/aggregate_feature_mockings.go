package mockings

import (
	"errors"
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/features/aggregate"
)

var ErrInvalidData = errors.New("invalid data")

type AggregateFeatureMocking struct{}

func (a *AggregateFeatureMocking) Aggregate(duration time.Duration, dimension string) (*aggregate.PostsStatAggregation, error) {
	return &aggregate.PostsStatAggregation{
		TotalPosts:       12,
		MinimumTimestamp: 1,
		MaximumTimestamp: 3,
		AvgLikes:         intP(2),
	}, nil
}

type AggregateFeatureErrorMocking struct{}

func (a *AggregateFeatureErrorMocking) Aggregate(duration time.Duration, dimension string) (*aggregate.PostsStatAggregation, error) {
	return nil, ErrInvalidData
}

func intP(i int) *int {
	return &i
}
