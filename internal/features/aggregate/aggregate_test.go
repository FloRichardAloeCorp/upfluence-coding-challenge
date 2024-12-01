package aggregate

import (
	"testing"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
)

func TestNewAggregateFeatures(t *testing.T) {
	type testData struct {
		name         string
		sseClient    *sse.SSEClient
		expectedRes1 AggregateFeatures
	}

	sseClient := &sse.SSEClient{}

	feature := NewAggregateFeatures(sseClient)

	if feature == nil {
		t.Error("aggregate feature factory creates a nil feature")
	}
}
