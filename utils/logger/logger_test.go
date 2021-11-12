package logger

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.EqualValues(t, envLogLevel, "LOG_LEVEL")
	assert.EqualValues(t, envLogOutput, "LOG_OUTPUT")
}

func TestGetLogger(t *testing.T) {
	myLogger := GetLogger()
	assert.NotNil(t, myLogger)
}

func TestGetLog(t *testing.T) {
	myLogger := GetLog()
	assert.NotNil(t, myLogger)
}

func TestInfo(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	initLogger(true)
	Info("my info message")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my info message")
}

func TestInfoWithField(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	initLogger(true)
	Info("my info message", Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my info message")
	assert.EqualValues(t, m["id"], "123")
}

func TestError(t *testing.T) {
	os.Setenv("LOG_LEVEL", "error")
	initLogger(true)
	Error("my error message", errors.New("new error"))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my error message")
	assert.EqualValues(t, m["error"], "new error")
}

func TestErrorWithField(t *testing.T) {
	os.Setenv("LOG_LEVEL", "error")
	initLogger(true)
	Error("my error message", errors.New("new error"), Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my error message")
	assert.EqualValues(t, m["error"], "new error")
	assert.EqualValues(t, m["id"], "123")
}

func TestDebug(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	initLogger(true)
	Debug("my debug message")
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my debug message")
}

func TestDebugWithField(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	initLogger(true)
	Debug("my debug message", Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my debug message")
	assert.EqualValues(t, m["id"], "123")
}

func TestPrint(t *testing.T) {
	initLogger(true)
	log.Print("a", "b")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "[a b]")
}

func TestPrintf(t *testing.T) {
	initLogger(true)
	log.Printf("my printf message")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my printf message")
}

func TestPrintfWithFormat(t *testing.T) {
	initLogger(true)
	log.Printf("my %s message", "formatted")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my formatted message")
}

func TestWarn(t *testing.T) {
	initLogger(true)
	Warn("my warn message")
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my warn message")
}

func TestWarnWithField(t *testing.T) {
	initLogger(true)
	Warn("my warn message", Field{
		Key:   "id",
		Value: "123",
	})
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "my warn message")
	assert.EqualValues(t, m["id"], "123")
}

func TestWriteInfo(t *testing.T) {
	initLogger(true)
	logMessage := "the is an info message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "info")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is an info message")
}

func TestWriteWarn(t *testing.T) {
	initLogger(true)
	logMessage := "the is a warning message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "warn")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is a warning message")
}

func TestWriteError(t *testing.T) {
	initLogger(true)
	logMessage := "the is an error message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "error")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is an error message")
}

func TestWriteDebug(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	initLogger(true)
	logMessage := "the is a debug message"
	written, writeErr := log.Write([]byte(logMessage))
	output := sink.String()
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(output), &m)
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, written)
	assert.Nil(t, writeErr)
	assert.EqualValues(t, len(logMessage), written)
	assert.EqualValues(t, m["level"], "debug")
	assert.Contains(t, m["caller"], "logger")
	assert.NotEmpty(t, m["time"])
	assert.EqualValues(t, m["msg"], "the is a debug message")
}
