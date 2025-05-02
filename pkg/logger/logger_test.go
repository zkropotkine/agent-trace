package logger

import (
	"context"
	"testing"

	testify "github.com/stretchr/testify/assert"
	"github.com/zkropotkine/agent-trace/config"
)

func TestNewLogrusLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		format   string
		expected string // substring expected in output
	}{
		{"Default JSON", "info", "json", ""},
		{"Text Format", "debug", "text", ""},
		{"Invalid Level", "notalevel", "text", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := testify.New(t)

			logger := NewLogRusLogger(config.Log{Level: tt.level, Format: tt.format})
			assert.NotNil(logger)
		})
	}
}

func TestContextLoggerInjection(t *testing.T) {
	ctx := context.Background()
	logger := NewLogRusLogger(config.Log{Level: "debug", Format: "json"})
	ctxWithLogger := WithLogger(ctx, logger)

	retrieved := FromContext(ctxWithLogger)

	assert := testify.New(t)
	assert.NotNil(retrieved)
	assert.Equal(logger, retrieved)
}

func TestContextLoggerFallback(t *testing.T) {
	nilCtx := context.Background()
	logger := FromContext(nilCtx)

	assert := testify.New(t)
	assert.NotNil(logger)
}

func TestLoggerWithFields(t *testing.T) {
	logger := NewLogRusLogger(config.Log{Level: "info", Format: "json"})
	logWithField := logger.WithField("testKey", "testValue")

	assert := testify.New(t)
	assert.NotNil(logWithField)

	logWithFields := logger.WithFields(map[string]interface{}{"key1": 1, "key2": true})
	assert.NotNil(logWithFields)
}
