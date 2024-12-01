package aggregate

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"time"
)

var (
	_                   AggregateFeatures = (*aggregateController)(nil)
	ErrNoPostsAvailable                   = errors.New("no posts available")
	ErrUnknownDimension                   = errors.New("unknown dimension")
)

type aggregateController struct {
	postStatsRepository iPostStatsRepository
}

func newAggregateController(postStatsRepository iPostStatsRepository) *aggregateController {
	return &aggregateController{
		postStatsRepository: postStatsRepository,
	}
}

func (c *aggregateController) Aggregate(duration time.Duration, dimension string) (*PostsStatAggregation, error) {
	poststats, err := c.postStatsRepository.ReadFor(duration)
	if err != nil {
		return nil, fmt.Errorf("can't read aggregate by id: %w", err)
	}

	if len(poststats) == 0 {
		return nil, ErrNoPostsAvailable
	}

	oldestPost := slices.MinFunc(poststats, func(a, b postStats) int {
		return cmp.Compare(a.Timestamp, b.Timestamp)
	})

	latestPost := slices.MaxFunc(poststats, func(a, b postStats) int {
		return cmp.Compare(a.Timestamp, b.Timestamp)
	})

	aggregation := &PostsStatAggregation{
		TotalPosts:       len(poststats),
		MinimumTimestamp: oldestPost.Timestamp,
		MaximumTimestamp: latestPost.Timestamp,
	}

	switch dimension {
	case "likes":
		aggregation.AvgLikes = intP(c.computeAvgLikes(poststats))
		return aggregation, nil
	case "comments":
		aggregation.AvgComments = intP(c.computeAvgComments(poststats))
		return aggregation, nil
	case "favorites":
		aggregation.AvgFavorites = intP(c.computeAvgFavorites(poststats))
		return aggregation, nil
	case "retweets":
		aggregation.AvgRetweets = intP(c.computeAvgRetweets(poststats))
		return aggregation, nil
	default:
		return nil, ErrUnknownDimension
	}
}

func (c *aggregateController) computeAvgLikes(postsStats []postStats) int {
	if len(postsStats) == 0 {
		return 0
	}

	sum := 0
	for _, stat := range postsStats {
		sum += stat.Likes
	}

	return sum / len(postsStats)
}

func (c *aggregateController) computeAvgComments(postsStats []postStats) int {
	if len(postsStats) == 0 {
		return 0
	}

	sum := 0
	for _, stat := range postsStats {
		sum += stat.Comments
	}

	return sum / len(postsStats)
}

func (c *aggregateController) computeAvgFavorites(postsStats []postStats) int {
	if len(postsStats) == 0 {
		return 0
	}

	sum := 0
	for _, stat := range postsStats {
		sum += stat.Favorites
	}

	return sum / len(postsStats)
}

func (c *aggregateController) computeAvgRetweets(postsStats []postStats) int {
	if len(postsStats) == 0 {
		return 0
	}

	sum := 0
	for _, stat := range postsStats {
		sum += stat.Retweets
	}

	return sum / len(postsStats)
}

func intP(i int) *int {
	return &i
}
