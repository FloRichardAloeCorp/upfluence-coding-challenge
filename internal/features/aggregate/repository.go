package aggregate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
)

var (
	_ iPostStatsRepository = (*postStatsRepository)(nil)

	ErrTooManyPosts     = errors.New("too many posts returned from stream")
	ErrEmptyEvent       = errors.New("empty event")
	ErrClosedSubscriber = errors.New("subscriber channel is closed")
)

type iPostStatsRepository interface {
	ReadFor(duration time.Duration) ([]postStats, error)
}

type postStatsRepository struct {
	sseClient *sse.SSEClient
}

func (r *postStatsRepository) ReadFor(duration time.Duration) ([]postStats, error) {
	sub, err := r.sseClient.NewSubscriber()
	if err != nil {
		return nil, fmt.Errorf("can't subscribe to sse server: %w", err)
	}
	defer r.sseClient.RemoveSubscriber(sub.ID)

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	postsStats := make([]postStats, 0)

	for {
		select {
		case event, ok := <-sub.Channel:
			if !ok {
				return nil, ErrClosedSubscriber
			}

			postStat, err := r.decodeEvent(event)
			if err != nil {
				return nil, fmt.Errorf("can't decode event: %w", err)
			}

			postsStats = append(postsStats, *postStat)
		case <-ctx.Done():
			return postsStats, nil
		}
	}
}

func (r *postStatsRepository) decodeEvent(event []byte) (*postStats, error) {
	rawPayload := make(map[string]json.RawMessage)
	if err := json.Unmarshal(event, &rawPayload); err != nil {
		return nil, fmt.Errorf("can't unmarshal event: %w", err)
	}

	// Event must contains only one entry
	if len(rawPayload) != 1 {
		return nil, ErrTooManyPosts
	}

	for _, postPayload := range rawPayload {
		postStat := postStats{}
		if err := json.Unmarshal(postPayload, &postStat); err != nil {
			return nil, fmt.Errorf("can't unmarshal post payload: %w", err)
		}

		return &postStat, nil
	}

	return nil, ErrEmptyEvent
}
