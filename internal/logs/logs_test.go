//nolint:goconst
package logs

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	// Success case with stderr
	config := Config{
		Level: "INFO",
	}

	log, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	if log.level != "INFO" {
		t.Fatalf("Expected %s level, got %s", "INFO", log.level)
	}

	if log.output != os.Stderr {
		t.Fatalf("Expected stderr output, got %s", log.output.Name())
	}

	// Success case with a file
	dir := t.TempDir()
	logFileName := dir + "/log.txt"
	_, err = os.Create(logFileName)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	config = Config{
		Level:      "INFO",
		OutputPath: logFileName,
	}

	log, err = NewLogger(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	if log.output.Name() != logFileName {
		t.Fatalf("Expected %s output file, got %s", logFileName, log.output.Name())
	}

	// Fail case: file doesn't exists
	config = Config{
		Level:      "INFO",
		OutputPath: "unknown",
	}

	_, err = NewLogger(config)
	if err == nil {
		t.Fatalf("An error is expected")
	}
}

func TestLoggerInfo(t *testing.T) {
	dir := t.TempDir()
	logFileName := dir + "/log.txt"
	_, err := os.Create(logFileName)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	config := Config{
		Level:      "INFO",
		OutputPath: logFileName,
	}

	log, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	log.Info("test", Field{Key: "test", Value: "value"})

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
		"test",
	}

	for _, expectedField := range expectedFields {
		if _, ok := logContent[expectedField]; !ok {
			t.Fatalf("expected field %s to be in log line but wasn't", expectedField)
		}
	}
}

func TestLoggerInfoWithErrorLevel(t *testing.T) {
	dir := t.TempDir()
	logFileName := dir + "/log.txt"
	_, err := os.Create(logFileName)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	config := Config{
		Level:      "ERROR",
		OutputPath: logFileName,
	}

	log, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	log.Info("test")

	content, err := os.ReadFile(logFileName)
	if err != nil {
		t.Fatalf("unable to read log file: %v ", err)
	}

	if len(content) != 0 {
		t.Fatalf("log file should be empty")
	}
}

func TestLoggerError(t *testing.T) {
	dir := t.TempDir()
	logFileName := dir + "/log.txt"
	_, err := os.Create(logFileName)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	config := Config{
		Level:      "ERROR",
		OutputPath: logFileName,
	}

	log, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	log.Error("test", Field{Key: "test", Value: "value"})

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
		"test",
	}

	for _, expectedField := range expectedFields {
		if _, ok := logContent[expectedField]; !ok {
			t.Fatalf("expected field %s to be in log line but wasn't", expectedField)
		}
	}
}

func TestLoggerErrorWithInfoLevel(t *testing.T) {
	dir := t.TempDir()
	logFileName := dir + "/log.txt"
	_, err := os.Create(logFileName)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	config := Config{
		Level:      "INFO",
		OutputPath: logFileName,
	}

	log, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v ", err)
	}

	log.Error("test", Field{Key: "test", Value: "value"})

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
		"test",
	}

	for _, expectedField := range expectedFields {
		if _, ok := logContent[expectedField]; !ok {
			t.Fatalf("expected field %s to be in log line but wasn't", expectedField)
		}
	}
}

func TestGetCallerInfo(t *testing.T) {
	type testData struct {
		name            string
		skip            int
		expectedPackage string
		expectedFile    string
	}

	testCases := [...]testData{
		{
			name:            "Success case: skip set to 1",
			skip:            1,
			expectedPackage: "logs",
			expectedFile:    "logs/logs_test.go",
		},
		{
			name:            "Success case: skip set to 100000",
			skip:            100000,
			expectedPackage: "unknown",
			expectedFile:    "unknown",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			packageName, file, line := getCallerInfo(testCase.skip)

			if testCase.expectedPackage != packageName {
				t.Errorf("Expected %s, got %s", testCase.expectedPackage, packageName)
			}

			if !strings.Contains(file, testCase.expectedFile) {
				t.Errorf("Expected %s to be in %s but isn't", file, testCase.expectedFile)
			}

			if line < 0 {
				t.Errorf("caller line shouldn't be lower thant 0")
			}
		})
	}
}
