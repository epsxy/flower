package log

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	cases := map[string]struct {
		input    LogLevel
		expected string
		isErr    bool
		errMsg   string
	}{
		"should convert log level `debug`": {
			input:    LogLevelDebug,
			expected: "debug",
		},
		"should convert log level `info`": {
			input:    LogLevelInfo,
			expected: "info",
		},
		"should convert log level `warn`": {
			input:    LogLevelWarn,
			expected: "warn",
		},
		"should convert log level `error`": {
			input:    LogLevelError,
			expected: "error",
		},
	}

	for _, c := range cases {
		require.Equal(t, c.input.String(), c.expected)
	}
}

func Test_Set(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected LogLevel
		isErr    bool
		errMsg   string
	}{
		"should convert log level `debug`": {
			input:    "debug",
			expected: LogLevelDebug,
		},
		"should convert log level `info`": {
			input:    "info",
			expected: LogLevelInfo,
		},
		"should convert log level `warn`": {
			input:    "warn",
			expected: LogLevelWarn,
		},
		"should convert log level `error`": {
			input:    "error",
			expected: LogLevelError,
		},
		"should fail while converting something else": {
			input:  "unsupported",
			isErr:  true,
			errMsg: "must be one of: `debug`, `info`, `warn` or `error`",
		},
	}

	for _, c := range cases {
		var logLevel LogLevel
		err := logLevel.Set(c.input)
		if c.isErr {
			require.NotNil(t, err, "error was expected but was nil")
			require.Contains(t, err.Error(), c.errMsg)
		} else {
			require.Equal(t, logLevel, c.expected)
		}
	}
}

func Test_ToSlogLevel(t *testing.T) {
	cases := map[string]struct {
		input    LogLevel
		expected slog.Level
	}{
		"should convert log level `debug`": {
			input:    LogLevelDebug,
			expected: slog.LevelDebug,
		},
		"should convert log level `info`": {
			input:    LogLevelInfo,
			expected: slog.LevelInfo,
		},
		"should convert log level `warn`": {
			input:    LogLevelWarn,
			expected: slog.LevelWarn,
		},
		"should convert log level `error`": {
			input:    LogLevelError,
			expected: slog.LevelError,
		},
	}

	for _, c := range cases {
		require.Equal(t, c.input.ToSlogLevel(), c.expected)
	}
}
