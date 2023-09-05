package zap

import (
	"github.com/lruggieri/fxnow/common/logger"
)

type Config struct {
	Development      bool
	Level            logger.Level
	OutputPaths      []string
	ErrorOutputPaths []string
	SkipLineEnding   bool
	AddCallerSkip    int
}
