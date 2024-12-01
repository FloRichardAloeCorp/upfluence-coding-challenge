package logs

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Logger represents a structured logging instance with configurable log levels
// and output destinations (e.g., stderr or a specified file).
type Logger struct {
	level  string
	output *os.File
}

type Config struct {
	// Level of the logger, can be one of INFO or ERROR.
	Level string `json:"level"`

	// OutputPath of the logger. Leave empty to write in stderr.
	OutputPath string `json:"output_path,omitempty"`
}

type log struct {
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
	Caller  string    `json:"caller"`
	Package string    `json:"package"`
	Msg     string    `json:"msg"`
}

func NewLogger(config Config) (*Logger, error) {
	var err error
	output := os.Stderr

	if config.OutputPath != "" && config.OutputPath != "stderr" {
		output, err = os.OpenFile(config.OutputPath, os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, fmt.Errorf("can't open log output file :%w", err)
		}
	}

	logger := &Logger{
		level:  config.Level,
		output: output,
	}

	return logger, nil
}

func (l *Logger) Info(msg string, fields ...Field) {
	if l.level != "INFO" {
		return
	}

	l.print(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.print(msg, fields...)
}

func (l *Logger) print(msg string, extraFields ...Field) {
	packageName, funcName, line := getCallerInfo(3)

	fields := []Field{
		{Key: "level", Value: l.level},
		{Key: "time", Value: time.Now().UTC().String()},
		{Key: "caller", Value: funcName + ":" + strconv.Itoa(line)},
		{Key: "package", Value: packageName},
		{Key: "msg", Value: msg},
	}

	fields = append(fields, extraFields...)

	content := encodeFieldsToJSON(fields...)
	fmt.Fprintln(l.output, content)
}

func getCallerInfo(skip int) (string, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown", 0
	}

	funcName := runtime.FuncForPC(pc).Name()

	parts := strings.Split(funcName, "/")
	funcParts := strings.Split(parts[len(parts)-1], ".")
	packageName := funcParts[0]

	return packageName, file, line
}
