package common

import (
	"log/slog"
	"os"
)

var Log = initLogger()

func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
