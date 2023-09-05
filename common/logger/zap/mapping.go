package zap

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"github.com/lruggieri/fxnow/common/logger"
)

func ToZapLevel(l logger.Level) zapcore.Level {
	switch l {
	case logger.LevelDebug:
		return zapcore.DebugLevel
	case logger.LevelInfo:
		return zapcore.InfoLevel
	case logger.LevelError:
		return zapcore.ErrorLevel
	default:
		panic(errors.Errorf("not support this level: %v", l))
	}
}

func FromZapLevel(l zapcore.Level) logger.Level {
	switch l {
	case zapcore.DebugLevel:
		return logger.LevelDebug
	case zapcore.InfoLevel:
		return logger.LevelInfo
	case zapcore.ErrorLevel:
		return logger.LevelError
	default:
		panic(errors.Errorf("not support this level: %v", l))
	}
}
