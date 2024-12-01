package logs

import (
	"testing"
)

func TestFieldToJSON(t *testing.T) {
	field := Field{
		Key:   "a",
		Value: "b",
	}

	if field.toJSON() != `"a":"b"` {
		t.Errorf("missmatching JSON field representation: have %s, want %s", field.toJSON(), `"a":"b"`)
	}
}
