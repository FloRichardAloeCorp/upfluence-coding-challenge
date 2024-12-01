package aggregate

import (
	"fmt"
	"testing"
	"time"
)

type postStatsRepositoryMocking struct {
	returnError bool
	NoResults   bool
}

func (r *postStatsRepositoryMocking) ReadFor(duration time.Duration) ([]postStats, error) {
	if r.returnError {
		return nil, fmt.Errorf("error")
	}

	if r.NoResults {
		return nil, nil
	}

	return []postStats{
		{
			Likes:     1,
			Comments:  2,
			Favorites: 3,
			Retweets:  4,
			Timestamp: 5,
		},
		{
			Likes:     7,
			Comments:  8,
			Favorites: 9,
			Retweets:  10,
			Timestamp: 11,
		},
	}, nil
}

func equalPostsStatAggregation(a, b PostsStatAggregation) bool {
	if a.TotalPosts != b.TotalPosts ||
		a.MinimumTimestamp != b.MinimumTimestamp ||
		a.MaximumTimestamp != b.MaximumTimestamp {
		return false
	}
	if (a.AvgLikes == nil) != (b.AvgLikes == nil) || (a.AvgLikes != nil && *a.AvgLikes != *b.AvgLikes) {
		return false
	}
	if (a.AvgComments == nil) != (b.AvgComments == nil) || (a.AvgComments != nil && *a.AvgComments != *b.AvgComments) {
		return false
	}
	if (a.AvgFavorites == nil) != (b.AvgFavorites == nil) || (a.AvgFavorites != nil && *a.AvgFavorites != *b.AvgFavorites) {
		return false
	}
	if (a.AvgRetweets == nil) != (b.AvgRetweets == nil) || (a.AvgRetweets != nil && *a.AvgRetweets != *b.AvgRetweets) {
		return false
	}
	return true
}

func TestNewAggregateController(t *testing.T) {
	repo := &postStatsRepositoryMocking{}
	controller := newAggregateController(repo)
	if controller.postStatsRepository != repo {
		t.Error("repository missmatch")
	}
}

func TestAggregateControllerAggregate(t *testing.T) {
	type testData struct {
		name           string
		shouldFail     bool
		mock           iPostStatsRepository
		duration       time.Duration
		dimension      string
		expectedResult *PostsStatAggregation
	}

	var testCases = [...]testData{
		{
			name:       "Success case with likes",
			shouldFail: false,
			mock: &postStatsRepositoryMocking{
				returnError: false,
				NoResults:   false,
			},
			duration:  time.Duration(5 * time.Second),
			dimension: "likes",
			expectedResult: &PostsStatAggregation{
				TotalPosts:       2,
				MinimumTimestamp: 5,
				MaximumTimestamp: 11,
				AvgLikes:         intP(4),
				AvgComments:      nil,
				AvgFavorites:     nil,
				AvgRetweets:      nil,
			},
		},
		{
			name:       "Success case with comments",
			shouldFail: false,
			mock: &postStatsRepositoryMocking{
				returnError: false,
				NoResults:   false,
			},
			duration:  time.Duration(5 * time.Second),
			dimension: "comments",
			expectedResult: &PostsStatAggregation{
				TotalPosts:       2,
				MinimumTimestamp: 5,
				MaximumTimestamp: 11,
				AvgLikes:         nil,
				AvgComments:      intP(5),
				AvgFavorites:     nil,
				AvgRetweets:      nil,
			},
		},
		{
			name:       "Success case with retweets",
			shouldFail: false,
			mock: &postStatsRepositoryMocking{
				returnError: false,
				NoResults:   false,
			},
			duration:  time.Duration(5 * time.Second),
			dimension: "retweets",
			expectedResult: &PostsStatAggregation{
				TotalPosts:       2,
				MinimumTimestamp: 5,
				MaximumTimestamp: 11,
				AvgLikes:         nil,
				AvgComments:      nil,
				AvgFavorites:     nil,
				AvgRetweets:      intP(7),
			},
		},
		{
			name:       "Success case with favorites",
			shouldFail: false,
			mock: &postStatsRepositoryMocking{
				returnError: false,
				NoResults:   false,
			},
			duration:  time.Duration(5 * time.Second),
			dimension: "favorites",
			expectedResult: &PostsStatAggregation{
				TotalPosts:       2,
				MinimumTimestamp: 5,
				MaximumTimestamp: 11,
				AvgLikes:         nil,
				AvgComments:      nil,
				AvgFavorites:     intP(6),
				AvgRetweets:      nil,
			},
		},
		{
			name:       "Fail case: repository returns an error",
			shouldFail: true,
			mock: &postStatsRepositoryMocking{
				returnError: true,
				NoResults:   false,
			},
			duration:  time.Duration(5 * time.Second),
			dimension: "likes",
		},
		{
			name:       "Fail case: repository returns 0 posts",
			shouldFail: true,
			mock: &postStatsRepositoryMocking{
				returnError: false,
				NoResults:   true,
			},
			duration:  time.Duration(5 * time.Second),
			dimension: "likes",
		},
		{
			name:       "Fail case: unknown dimension",
			shouldFail: true,
			mock: &postStatsRepositoryMocking{
				returnError: false,
				NoResults:   false,
			},
			duration:  time.Duration(5 * time.Second),
			dimension: "invalid",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := &aggregateController{
				postStatsRepository: testCase.mock,
			}
			stats, err := instance.Aggregate(testCase.duration, testCase.dimension)
			if testCase.shouldFail {
				if err == nil {
					t.Errorf("expected an error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !equalPostsStatAggregation(*stats, *testCase.expectedResult) {
					t.Errorf("expected %v, got %v", *testCase.expectedResult, *stats)
				}
			}
		})
	}
}

func TestAggregateControllerComputeAvgLikes(t *testing.T) {
	type testData struct {
		name           string
		postsStats     []postStats
		expectedResult int
	}

	var testCases = [...]testData{
		{
			name: "Succes case",
			postsStats: []postStats{
				{
					Likes: 2,
				},
				{
					Likes: 2,
				},
			},
			expectedResult: 2,
		},
		{
			name:           "Succes case: empty posts",
			postsStats:     []postStats{},
			expectedResult: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := &aggregateController{
				postStatsRepository: &postStatsRepositoryMocking{
					returnError: false,
					NoResults:   false,
				},
			}

			avgLikes := instance.computeAvgLikes(testCase.postsStats)
			if avgLikes != testCase.expectedResult {
				t.Errorf("expected %v, got %v", avgLikes, testCase.expectedResult)
			}
		})
	}
}

func TestAggregateControllerComputeAvgComments(t *testing.T) {
	type testData struct {
		name           string
		postsStats     []postStats
		expectedResult int
	}

	var testCases = [...]testData{
		{
			name: "Succes case",
			postsStats: []postStats{
				{
					Comments: 2,
				},
				{
					Comments: 2,
				},
			},
			expectedResult: 2,
		},
		{
			name:           "Succes case: empty posts",
			postsStats:     []postStats{},
			expectedResult: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := &aggregateController{
				postStatsRepository: &postStatsRepositoryMocking{
					returnError: false,
					NoResults:   false,
				},
			}

			avgComments := instance.computeAvgComments(testCase.postsStats)
			if avgComments != testCase.expectedResult {
				t.Errorf("expected %v, got %v", avgComments, testCase.expectedResult)
			}
		})
	}
}

func TestAggregateControllerComputeAvgFavorites(t *testing.T) {
	type testData struct {
		name           string
		postsStats     []postStats
		expectedResult int
	}

	var testCases = [...]testData{
		{
			name: "Succes case",
			postsStats: []postStats{
				{
					Favorites: 2,
				},
				{
					Favorites: 2,
				},
			},
			expectedResult: 2,
		},
		{
			name:           "Succes case: empty posts",
			postsStats:     []postStats{},
			expectedResult: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := &aggregateController{
				postStatsRepository: &postStatsRepositoryMocking{
					returnError: false,
					NoResults:   false,
				},
			}

			avgFavorites := instance.computeAvgFavorites(testCase.postsStats)
			if avgFavorites != testCase.expectedResult {
				t.Errorf("expected %v, got %v", avgFavorites, testCase.expectedResult)
			}
		})
	}
}

func TestAggregateControllerComputeAvgRetweets(t *testing.T) {
	type testData struct {
		name           string
		postsStats     []postStats
		expectedResult int
	}

	var testCases = [...]testData{
		{
			name: "Succes case",
			postsStats: []postStats{
				{
					Retweets: 2,
				},
				{
					Retweets: 2,
				},
			},
			expectedResult: 2,
		},
		{
			name:           "Succes case: empty posts",
			postsStats:     []postStats{},
			expectedResult: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := &aggregateController{
				postStatsRepository: &postStatsRepositoryMocking{
					returnError: false,
					NoResults:   false,
				},
			}

			avgRetweets := instance.computeAvgRetweets(testCase.postsStats)
			if avgRetweets != testCase.expectedResult {
				t.Errorf("expected %v, got %v", avgRetweets, testCase.expectedResult)
			}
		})
	}
}

func TestIntP(t *testing.T) {
	i := 2

	iP := intP(i)

	if i != *iP {
		t.Errorf("expected %v, got %v", i, *iP)
	}
}
