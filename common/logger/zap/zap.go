package zap

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lruggieri/fxnow/common/logger"
)

const (
	ZapFieldKeyMessage = "msg"
	ZapFieldKeyTime    = "time"
	ZapFieldKeyError   = "error"
)

type ImplLogger struct {
	cfg       Config
	zapLogger *zap.Logger
	ctx       context.Context
}

func (z *ImplLogger) Config() Config {
	return z.cfg
}

// SetLevel implements logger.Logger
func (z *ImplLogger) SetLevel(level logger.Level) logger.Logger {
	cfg := z.cfg
	cfg.Level = level

	newInstance := New(cfg)

	return z.copyWith(
		z.zapLogger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return newInstance.zapLogger.Core()
		})),
		z.ctx,
	)
}

func (z *ImplLogger) WithField(key string, value interface{}) logger.Logger {
	return z.copyWith(
		z.zapLogger.With(zap.Any(key, value)),
		z.ctx,
	)
}

func (z *ImplLogger) WithFields(fields logger.Fields) logger.Logger {
	var (
		zapFields = make([]zap.Field, len(fields))
		i         = 0
	)

	for k, v := range fields {
		zapFields[i] = zap.Any(k, v)
		i++
	}

	return z.copyWith(
		z.zapLogger.With(zapFields...),
		z.ctx,
	)
}

func (z *ImplLogger) WithError(err error) logger.Logger {
	return z.WithField(ZapFieldKeyError, fmt.Sprintf("%+v", err))
}

func (z *ImplLogger) WithContext(ctx context.Context) logger.Logger {
	return z.copyWith(
		z.zapLogger,
		ctx,
	)
}

func (z *ImplLogger) Debug(format string, args ...interface{}) {
	z.write(zapcore.DebugLevel, format, args...)
}

func (z *ImplLogger) Info(format string, args ...interface{}) {
	z.write(zapcore.InfoLevel, format, args...)
}

func (z *ImplLogger) Error(format string, args ...interface{}) {
	z.write(zapcore.ErrorLevel, format, args...)
}

func (z *ImplLogger) write(level zapcore.Level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ce := z.zapLogger.Check(level, msg); ce != nil {
		fields := []zapcore.Field{}

		// TODO enricher logic

		ce.Write(fields...)
	}
}

func (z *ImplLogger) Context() context.Context {
	return z.ctx
}

//nolint:revive
func (z *ImplLogger) copyWith(
	zapLogger *zap.Logger,
	ctx context.Context,
) *ImplLogger {
	return &ImplLogger{
		zapLogger: zapLogger,
		ctx:       ctx,
	}
}

func New(cfg Config) *ImplLogger {
	var (
		zapConfig zap.Config
		instance  = &ImplLogger{
			ctx: context.Background(),
			cfg: cfg,
		}
	)

	if cfg.Development {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	if len(cfg.OutputPaths) > 0 {
		zapConfig.OutputPaths = cfg.OutputPaths
	}

	if len(cfg.ErrorOutputPaths) > 0 {
		zapConfig.ErrorOutputPaths = cfg.ErrorOutputPaths
	}

	zapConfig.EncoderConfig.SkipLineEnding = cfg.SkipLineEnding // avoid printing a newline at the end
	zapConfig.EncoderConfig.MessageKey = ZapFieldKeyMessage
	zapConfig.EncoderConfig.TimeKey = ZapFieldKeyTime
	zapConfig.Level.SetLevel(ToZapLevel(cfg.Level))

	if cfg.AddCallerSkip == 0 {
		cfg.AddCallerSkip = 3
	}

	instance.zapLogger, _ = zapConfig.Build(
		zap.AddCallerSkip(cfg.AddCallerSkip),
	)

	return instance
}

func NewTestLogger() logger.Logger {
	return New(Config{
		Development:    true,
		Level:          logger.LevelDebug,
		SkipLineEnding: false,
	})
}
