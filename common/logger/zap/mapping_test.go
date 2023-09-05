package zap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/lruggieri/fxnow/common/logger"
)

func TestToZapLevel(t *testing.T) {
	assert.Equal(t, zapcore.DebugLevel, ToZapLevel(logger.LevelDebug))
	assert.Equal(t, zapcore.InfoLevel, ToZapLevel(logger.LevelInfo))
	assert.Equal(t, zapcore.ErrorLevel, ToZapLevel(logger.LevelError))
	assert.Panics(t, func() {
		ToZapLevel(logger.Level(4))
	})
}

func TestFromZapLevel(t *testing.T) {
	assert.Equal(t, logger.LevelDebug, FromZapLevel(zapcore.DebugLevel))
	assert.Equal(t, logger.LevelInfo, FromZapLevel(zapcore.InfoLevel))
	assert.Equal(t, logger.LevelError, FromZapLevel(zapcore.ErrorLevel))

	assert.Panics(t, func() {
		FromZapLevel(zapcore.WarnLevel)
	})
	assert.Panics(t, func() {
		FromZapLevel(zapcore.PanicLevel)
	})
	assert.Panics(t, func() {
		FromZapLevel(zapcore.FatalLevel)
	})
}
