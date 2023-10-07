package global

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

var (
	verbose = false
	dryRun  = false
	logger  *slog.Logger
)

func SetVerbose(v bool) {
	verbose = v
}

func GetVerbose() bool {
	return verbose
}

func SetDryRun(d bool) {
	dryRun = d
}

func GetDryRun() bool {
	return dryRun
}

func SetLogger(level slog.Level) {
	w := os.Stderr

	// create a new logger
	logger = slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      level,
			TimeFormat: time.Kitchen,
		}),
	)
}

func GetLogger() *slog.Logger {
	return logger
}
