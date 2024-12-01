package middlewares

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/FloRichardAloeCorp/upfluence-coding-challenge/internal/logs"
	"github.com/gin-gonic/gin"
)

func TestRequestsLogger(t *testing.T) {
	dir := t.TempDir()
	logFileName := dir + "/log.txt"
	os.Create(logFileName)

	config := logs.Config{
		Level:      "INFO",
		OutputPath: logFileName,
	}

	log, err := logs.NewLogger(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = httptest.NewRequest("GET", "/analysis", nil)

	RequestsLogger(log)(ctx)

	content, err := os.ReadFile(logFileName)
	if err != nil {
		t.Fatalf("unable to read log file: %v ", err)
	}

	if len(content) == 0 {
		t.Fatalf("log file should not be empty")
	}

	logContent := make(map[string]string)
	if err := json.Unmarshal(content, &logContent); err != nil {
		t.Fatalf("log line is not unmarshalable to json: %v", err)
	}

	expectedFields := []string{
		"level",
		"time",
		"caller",
		"package",
		"msg",
		"path",
		"method",
		"latency",
	}

	for _, expectedField := range expectedFields {
		if _, ok := logContent[expectedField]; !ok {
			t.Fatalf("expected field %s to be in log line but wasn't", expectedField)
		}
	}
}
