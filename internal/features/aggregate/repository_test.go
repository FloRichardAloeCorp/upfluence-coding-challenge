package aggregate

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
)

var (
	loggerInstance, _ = logs.NewLogger(logs.Config{
		Level: "INFO",
	})

	eventData = `data: {"tweet":{"id":959084760,"content":"Wishing for the heat of summer ðŸ”¥ https://t.co/Ykb72ulGdR","retweets":19,"favorites":643,"timestamp":1681859460,"post_id":"1648464174270521347","is_retweet":false,"comments":24}}`
)

func createSSEServerMock(interval time.Duration, data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			w.WriteHeader(500)
			return
		}

		for i := 0; i < 5; i++ {
			_, _ = w.Write(append(data, []byte("\n\n")...))
			flusher.Flush()
			time.Sleep(interval)
		}
	}))
}

func TestPostStatsRepositoryReadFor(t *testing.T) {
	server := createSSEServerMock(1*time.Second, []byte(eventData))
	defer server.Close()

	sseClientConfig := sse.Config{
		ServerURL:               server.URL,
		MaxReconnectionAttempts: 1,
	}

	sseClient := sse.NewSSEClient(sseClientConfig, loggerInstance)

	var listenErr error
	go func() {
		listenErr = sseClient.Listen()
	}()
	defer sseClient.Close()

	repo := postStatsRepository{
		sseClient: sseClient,
	}

	posts, err := repo.ReadFor(4 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error, got %v", err)
	}

	if listenErr != nil {
		t.Fatalf("unexpected error from sseClient.Listen, got %v", err)
	}

	if len(posts) == 0 {
		t.Errorf("post stats should not be empty")
	}
}

func TestPostStatsRepositoryReadForInvalidEvent(t *testing.T) {
	server := createSSEServerMock(1*time.Second, []byte("data: invalid"))
	defer server.Close()

	sseClientConfig := sse.Config{
		ServerURL:               server.URL,
		MaxReconnectionAttempts: 1,
	}

	sseClient := sse.NewSSEClient(sseClientConfig, loggerInstance)
	var listenErr error
	go func() {
		listenErr = sseClient.Listen()
	}()
	defer sseClient.Close()

	repo := postStatsRepository{
		sseClient: sseClient,
	}

	posts, err := repo.ReadFor(4 * time.Second)
	if err == nil {
		t.Fatalf("expected error %v, got %v", ErrClosedSubscriber, err)
	}

	if listenErr != nil {
		t.Fatalf("unexpected error from sseClient.Listen, got %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("post stats should be empty")
	}
}

func TestPostStatsRepositoryDecodeEvent(t *testing.T) {
	type testData struct {
		name           string
		event          []byte
		shouldFail     bool
		expectedResult *postStats
	}

	testCases := [...]testData{
		{
			name:       "Success case",
			event:      []byte(`{"yt":{"likes":2,"timestamp":1}}`),
			shouldFail: false,
			expectedResult: &postStats{
				Likes:     2,
				Comments:  0,
				Favorites: 0,
				Retweets:  0,
				Timestamp: 1,
			},
		},
		{
			name:       "Fail case: event is not a json string",
			event:      []byte(`invalid`),
			shouldFail: true,
		},
		{
			name:       "Fail case: multiple key",
			event:      []byte(`{"a":{"likes":2,"timestamp":1},"b":{"likes":2,"timestamp":1}}`),
			shouldFail: true,
		},
		{
			name:       "Fail case: invalid envent body",
			event:      []byte(`{"yt":{"likes":true}}`),
			shouldFail: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance := &postStatsRepository{}
			post, err := instance.decodeEvent(testCase.event)

			if testCase.shouldFail {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error")
				}

				if *post != *testCase.expectedResult {
					t.Errorf("expected %v got %v", post, testCase.expectedResult)
				}
			}
		})
	}
}
