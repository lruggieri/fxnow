package logger

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
)

type Level uint8

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	default:
		panic("invalid log level")
	}
}

func (l *Level) UnmarshalText(data []byte) (err error) {
	raw := string(data)

	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(error); ok {
				err = e
			} else {
				err = errors.Errorf("%v", p)
			}
		}
	}()

	if parsedLevel, err1 := l.tryParseLogLevelByNumber(raw); err1 == nil {
		*l = parsedLevel
	} else {
		*l = GetLoggingLevelFromString(raw)
	}

	return
}

func (*Level) tryParseLogLevelByNumber(raw string) (Level, error) {
	num, err1 := strconv.ParseInt(raw, 10, 8)

	if err1 == nil {
		if num < int64(LevelDebug) || num > int64(LevelError) {
			return 0, errors.Errorf("invalid log level: %d", num)
		}

		return Level(num), nil
	}

	return 0, err1
}

func GetLoggingLevelFromString(level string) Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "error":
		return LevelError
	default:
		panic(errors.Errorf("invalid log level: %s", level))
	}
}
