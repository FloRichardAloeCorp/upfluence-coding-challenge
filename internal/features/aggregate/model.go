package aggregate

type postStats struct {
	Likes     int   `json:"likes,omitempty"`
	Comments  int   `json:"comments,omitempty"`
	Favorites int   `json:"favorites,omitempty"`
	Retweets  int   `json:"retweets,omitempty"`
	Timestamp int64 `json:"timestamp"`
}

type PostsStatAggregation struct {
	TotalPosts       int   `json:"total_posts"`
	MinimumTimestamp int64 `json:"minimum_timestamp"`
	MaximumTimestamp int64 `json:"maximum_timestamp"`

	AvgLikes     *int `json:"avg_likes,omitempty"`
	AvgComments  *int `json:"avg_comments,omitempty"`
	AvgFavorites *int `json:"avg_favorites,omitempty"`
	AvgRetweets  *int `json:"avg_retweets,omitempty"`
}
