package zap_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lruggieri/fxnow/common/logger"
	"github.com/lruggieri/fxnow/common/logger/zap"
)

func TestNew(t *testing.T) {
	a := zap.New(zap.Config{})
	assert.NotNil(t, a)

	a = zap.New(zap.Config{
		Development: true,
	})
	assert.NotNil(t, a)

	a = zap.New(zap.Config{
		Development:      true,
		Level:            logger.LevelError,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stdout"},
		SkipLineEnding:   true,
	})
	assert.NotNil(t, a)

	assert.NotNil(t, zap.NewTestLogger())
}

func TestLogger_LogFunctions(t *testing.T) {
	sinkName := "TestLogger_LogFunctions"
	l := createLogger(sinkName)
	s := getSink(sinkName)

	l.Debug("debug %s", "msg")
	l.Info("info %s", "msg")
	l.Error("error %s", "msg")

	assert.Equal(t, "debug msg", s.line[0][zap.ZapFieldKeyMessage])
	assert.Equal(t, "DEBUG", s.line[0]["level"])

	assert.Equal(t, "info msg", s.line[1][zap.ZapFieldKeyMessage])
	assert.Equal(t, "INFO", s.line[1]["level"])

	assert.Equal(t, "error msg", s.line[2][zap.ZapFieldKeyMessage])
	assert.Equal(t, "ERROR", s.line[2]["level"])
}

func TestLogger_WithField(t *testing.T) {
	field1Key := "field1"
	field2Key := "field2"
	sinkName := "TestLogger_WithField"

	type Person struct {
		Name string
		Age  int
	}

	l := createLogger(sinkName)
	s := getSink(sinkName)

	l.
		WithField(field1Key, "field 1 value").
		WithField(field2Key, "field 2 value").
		WithField("Person", Person{Name: "John Doe", Age: 30}).
		Info("info msg")

	assert.Equal(t, "info msg", s.line[0][zap.ZapFieldKeyMessage])
	assert.Equal(t, "INFO", s.line[0]["level"])
	assert.Equal(t, "field 1 value", s.line[0][field1Key])
	assert.Equal(t, "field 2 value", s.line[0][field2Key])
	assert.Equal(t, map[string]interface{}{"Age": float64(30), "Name": "John Doe"}, s.line[0]["Person"])
}

func TestLogger_WithError(t *testing.T) {
	err := errors.New("test error")
	sinkName := "TestLogger_WithError"
	l := createLogger(sinkName)
	s := getSink(sinkName)

	l.WithError(err).Info("info msg")

	assert.Equal(t, "info msg", s.line[0][zap.ZapFieldKeyMessage])
	assert.Equal(t, "INFO", s.line[0]["level"])
	assert.Equal(t, err.Error(), s.line[0]["error"])
}
