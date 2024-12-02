package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"testing"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/features/aggregate"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/test/mockings"
	"github.com/gin-gonic/gin"
)

var (
	testConfig = AnalysisHandlerConfig{
		AuthorizedDimensions: []string{
			"likes",
		},
	}

	loggerInstance, _ = logs.NewLogger(logs.Config{
		Level: "INFO",
	})
)

func TestNewAnalysisHandler(t *testing.T) {
	feature := mockings.AggregateFeatureMocking{}

	handler := NewAnalysisHandler(testConfig, &feature, loggerInstance)

	if !slices.Equal(testConfig.AuthorizedDimensions, handler.authorizedDimension) {
		t.Errorf("AnalysisHandler authorized dimensions differ from the injected ones.")
	}

	if handler.aggregateFeatures != &feature {
		t.Errorf("AnalysisHandler aggregate feature differ from the injected one.")
	}
}

func TestAnalysisHandlerRegisterRoutes(t *testing.T) {
	router := gin.Default()

	handler := NewAnalysisHandler(testConfig, &mockings.AggregateFeatureMocking{}, loggerInstance)

	handler.RegisterRoutes(router)

	route := router.Routes()[0]
	if route.Path != "/analysis" && route.Method != "GET" {
		t.Errorf("Handler routes should be GET /analysis, got %s %s", route.Method, route.Path)
	}
}

func TestAnalysisHandlerGet(t *testing.T) {
	type testData struct {
		name                string
		queryParams         map[string]string
		authorizedDimension []string
		expectedStatusCode  int
		hasResponseBody     bool
	}

	testCases := [...]testData{
		{
			name: "Success case",
			queryParams: map[string]string{
				"duration":  "5s",
				"dimension": "likes",
			},
			authorizedDimension: []string{
				"likes",
			},
			expectedStatusCode: http.StatusOK,
			hasResponseBody:    true,
		},
		{
			name: "Fail case: no duration query param",
			queryParams: map[string]string{
				"dimension": "likes",
			},
			authorizedDimension: []string{
				"likes",
			},
			expectedStatusCode: http.StatusBadRequest,
			hasResponseBody:    false,
		},
		{
			name: "Fail case: duration query param is not a go duration",
			queryParams: map[string]string{
				"duration":  "invalid",
				"dimension": "likes",
			},
			authorizedDimension: []string{
				"likes",
			},
			expectedStatusCode: http.StatusBadRequest,
			hasResponseBody:    false,
		},
		{
			name: "Fail case: negative duration",
			queryParams: map[string]string{
				"duration":  "-5s",
				"dimension": "likes",
			},
			authorizedDimension: []string{
				"likes",
			},
			expectedStatusCode: http.StatusBadRequest,
			hasResponseBody:    false,
		},
		{
			name: "Fail case: no dimension query param",
			queryParams: map[string]string{
				"duration": "5s",
			},
			authorizedDimension: []string{
				"likes",
			},
			expectedStatusCode: http.StatusBadRequest,
			hasResponseBody:    false,
		},
		{
			name: "Fail case: invalid dimension value supplied",
			queryParams: map[string]string{
				"duration":  "5s",
				"dimension": "unknown",
			},
			authorizedDimension: []string{
				"likes",
			},
			expectedStatusCode: http.StatusBadRequest,
			hasResponseBody:    false,
		},
		{
			name: "Fail case: empty authorized dimensions blocks everything",
			queryParams: map[string]string{
				"duration":  "5s",
				"dimension": "likes",
			},
			authorizedDimension: []string{},
			expectedStatusCode:  http.StatusBadRequest,
			hasResponseBody:     false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Preparing gin context
			writer := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(writer)

			ctx.Request = httptest.NewRequest("GET", "/analysis", nil)

			values := url.Values{}
			for k, v := range testCase.queryParams {
				values[k] = []string{v}
			}
			ctx.Request.URL.RawQuery = values.Encode()

			instance := &AnalysisHandler{
				aggregateFeatures:   &mockings.AggregateFeatureMocking{},
				authorizedDimension: testCase.authorizedDimension,
				log:                 loggerInstance,
			}

			instance.Get(ctx)

			if writer.Code != testCase.expectedStatusCode {
				t.Fatalf("expected status code %d, got %d", testCase.expectedStatusCode, writer.Code)
			}

			if testCase.hasResponseBody {
				aggregation := aggregate.PostsStatAggregation{}
				responseBody, err := io.ReadAll(writer.Body)
				if err != nil {
					t.Errorf("should be able to read response body, error %v", err)
				}

				err = json.Unmarshal(responseBody, &aggregation)
				if err != nil {
					t.Errorf("should be able to unmarshal response body in a PostsStatAggregation, error %v", err)
				}
			}
		})
	}
}

func TestAnalysisHandlerGetAggregateFeatureError(t *testing.T) {
	// Preparing gin context
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)

	ctx.Request = httptest.NewRequest("GET", "/analysis", nil)

	values := url.Values{
		"duration":  []string{"5s"},
		"dimension": []string{"likes"},
	}

	ctx.Request.URL.RawQuery = values.Encode()

	instance := &AnalysisHandler{
		aggregateFeatures: &mockings.AggregateFeatureErrorMocking{},
		authorizedDimension: []string{
			"likes",
		},
		log: loggerInstance,
	}

	instance.Get(ctx)

	if writer.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, writer.Code)
	}
}
