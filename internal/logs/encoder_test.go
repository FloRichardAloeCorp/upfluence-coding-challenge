package logs

import (
	"testing"
)

func TestEncodeFieldsToJson(t *testing.T) {
	fields := []Field{
		{
			Key:   "a",
			Value: "b",
		},
		{
			Key:   "c",
			Value: "d",
		},
	}

	jsonString := encodeFieldsToJSON(fields...)

	if jsonString != `{"a":"b","c":"d"}` {
		t.Errorf("missmatching JSON fields representation: have %s, want %s", jsonString, `{"a":"b","c":"d"}`)
	}
}
