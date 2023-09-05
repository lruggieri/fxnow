package logger

import (
	"context"
	"errors"
	"sync"
)

var (
	initOnce       sync.Once
	loggerInstance Logger
)

// InitLogger : initialize common logger with the input one.
// Assumption: we don't need to use different type of loggers in different part of our codebase, therefore having
// a general logger reduces dependency injections and makes it easy to call for the logger
func InitLogger(logger Logger) {
	initOnce.Do(func() {
		loggerInstance = logger
	})
}

func Get() Logger {
	if loggerInstance == nil {
		panic(errors.New("should init logger before using"))
	}

	return loggerInstance
}

func SetLevel(level Level) Logger {
	return Get().SetLevel(level)
}

func WithField(key string, value interface{}) Logger {
	return Get().WithField(key, value)
}

func WithFields(fields Fields) Logger {
	return Get().WithFields(fields)
}

func WithContext(ctx context.Context) Logger {
	return Get().WithContext(ctx)
}

func WithError(err error) Logger {
	return Get().WithError(err)
}

func Debug(format string, args ...interface{}) {
	Get().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	Get().Info(format, args...)
}

func Error(format string, args ...interface{}) {
	Get().Error(format, args...)
}
