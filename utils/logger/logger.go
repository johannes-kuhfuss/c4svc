package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

var (
	log logger
)

type loggerInterface interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type logger struct {
	log *zap.Logger
}

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "msg"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.StacktraceKey = ""

	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{getOutput()}
	logConfig.Level = zap.NewAtomicLevelAt(getLevel())
	logConfig.Encoding = "json"
	logConfig.EncoderConfig = encoderConfig

	var err error
	if log.log, err = logConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1)); err != nil {
		panic(err)
	}
}

func getLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	output := strings.TrimSpace(os.Getenv(envLogOutput))
	if output == "" {
		return "stdout"
	}
	return output
}

func GetLogger() loggerInterface {
	return log
}

func GetLog() *zap.Logger {
	return log.log
}

func (l logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v...))
	}
}

func (l logger) Print(v ...interface{}) {
	Info(fmt.Sprintf("%v", v))
}

func Debug(msg string, tags ...zap.Field) {
	log.log.Debug(msg, tags...)
	log.log.Sync()
}

func Info(msg string, tags ...zap.Field) {
	log.log.Info(msg, tags...)
	log.log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.log.Error(msg, tags...)
	log.log.Sync()
}
