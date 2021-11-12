package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func setupLogsCapture() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)

	return zap.New(core), logs
}

func TestInfo(t *testing.T) {
	logger, logs := setupLogsCapture()
	logger.Info("my message")
	assert.NotNil(t, logs)
	assert.EqualValues(t, 1, logs.Len())
	assert.EqualValues(t, zap.InfoLevel, logs.All()[0].Level)
	assert.EqualValues(t, "my message", logs.All()[0].Message)
}
