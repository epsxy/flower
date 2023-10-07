package log

import (
	"errors"
	"log/slog"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *LogLevel) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *LogLevel) Set(v string) error {
	switch v {
	case "debug", "info", "warn", "error":
		*e = LogLevel(v)
		return nil
	default:
		return errors.New("must be one of: `debug`, `info`, `warn` or `error`")
	}
}

// Type is only used in help text
func (e *LogLevel) Type() string {
	return "logLevel"
}

func (e *LogLevel) ToSlogLevel() slog.Level {
	if *e == LogLevelDebug {
		return slog.LevelDebug
	}
	if *e == LogLevelInfo {
		return slog.LevelInfo
	}
	if *e == LogLevelWarn {
		return slog.LevelWarn
	}
	if *e == LogLevelError {
		return slog.LevelError
	}
	return slog.LevelError
}
