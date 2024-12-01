package aggregate

import (
	"cmp"
	"fmt"
	"slices"
	"time"
)

var (
	_ AggregateFeatures = (*aggregateController)(nil)
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
		return nil, fmt.Errorf("no stats available")
	}

	olderPost := slices.MinFunc(poststats, func(a, b postStats) int {
		return cmp.Compare(a.Timestamp, b.Timestamp)
	})

	latestPost := slices.MaxFunc(poststats, func(a, b postStats) int {
		return cmp.Compare(a.Timestamp, b.Timestamp)
	})

	aggregation := &PostsStatAggregation{
		TotalPosts:       len(poststats),
		MinimumTimestamp: olderPost.Timestamp,
		MaximumTimestamp: latestPost.Timestamp,
	}

	if dimension == "likes" {
		aggregation.AvgLikes = intP(c.computeAvgLikes(poststats))
		return aggregation, nil
	}

	if dimension == "comments" {
		aggregation.AvgComments = intP(c.computeAvgComments(poststats))
		return aggregation, nil
	}

	if dimension == "favorites" {
		aggregation.AvgFavorites = intP(c.computeAvgFavorites(poststats))
		return aggregation, nil
	}

	if dimension == "retweets" {
		aggregation.AvgRetweets = intP(c.computeAvgRetweets(poststats))
		return aggregation, nil
	}

	return aggregation, nil
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
