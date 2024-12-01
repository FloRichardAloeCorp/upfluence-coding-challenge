package aggregate

import (
	"testing"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/interfaces/sse"
)

func TestNewAggregateFeatures(t *testing.T) {
	sseClient := &sse.Client{}

	feature := NewAggregateFeatures(sseClient)

	if feature == nil {
		t.Error("aggregate feature factory creates a nil feature")
	}
}
