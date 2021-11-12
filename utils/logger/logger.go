package logger

import (
	"bytes"
	"fmt"
	"net/url"
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
	log  logger
	sink *MemorySink
)

type loggerInterface interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Write([]byte) (int, error)
}

type logger struct {
	log *zap.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

type MemorySink struct {
	*bytes.Buffer
}

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

func init() {
	initLogger(false)
}

func initLogger(test bool) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "msg"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.StacktraceKey = ""

	logConfig := zap.NewProductionConfig()
	if test {
		zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
			return sink, nil
		})
		sink = &MemorySink{new(bytes.Buffer)}
		logConfig.OutputPaths = []string{"memory://"}
	} else {
		logConfig.OutputPaths = []string{getOutput()}
	}

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

func (l logger) Write(data []byte) (n int, err error) {
	//fmt.Printf("data: %v", string(data))
	logMessage := string(data)
	if strings.Contains(strings.ToLower(logMessage), "warn") {
		Warn(logMessage)
	} else if strings.Contains(logMessage, "debug") {
		Debug(logMessage)
	} else {
		Info(logMessage)
	}

	return len(data), nil
}

func fieldsToZapField(tags []Field) []zapcore.Field {
	zapTags := make([]zap.Field, 0)
	for _, tag := range tags {
		zapTag := zap.Any(tag.Key, tag.Value)
		zapTags = append(zapTags, zapTag)
	}
	return zapTags
}

func Debug(msg string, tags ...Field) {
	zapTags := fieldsToZapField(tags)
	log.log.Debug(msg, zapTags...)
	log.log.Sync()
}

func Info(msg string, tags ...Field) {
	zapTags := fieldsToZapField(tags)
	log.log.Info(msg, zapTags...)
	log.log.Sync()
}

func Warn(msg string, tags ...Field) {
	zapTags := fieldsToZapField(tags)
	log.log.Warn(msg, zapTags...)
	log.log.Sync()
}

func Error(msg string, err error, tags ...Field) {
	zapTags := fieldsToZapField(tags)
	zapTags = append(zapTags, zap.NamedError("error", err))
	log.log.Error(msg, zapTags...)
	log.log.Sync()
}
