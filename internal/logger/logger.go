// Package logger provides a global logger instance for the application.
package logger

import (
	"log/slog"
	"os"

	"golang.org/x/term"

	"github.com/AtifChy/aiub-notice/internal/logger/handler"
)

var (
	logger *slog.Logger
	opts   slog.HandlerOptions
)

// L returns the global logger instance.
func L() *slog.Logger { return logger }

func Default() *slog.Logger {
	var rootHandler slog.Handler
	if isInteractive() {
		rootHandler = handler.NewPrettyHandler(os.Stderr, &opts)
	} else {
		rootHandler = slog.NewJSONHandler(os.Stderr, &opts)
	}

	return slog.New(rootHandler)
}

func SetOutputFile(logfile *os.File) {
	var rootHandler slog.Handler

	jsonHandler := slog.NewJSONHandler(logfile, &opts)
	if isInteractive() {
		prettyHandler := handler.NewPrettyHandler(os.Stderr, &opts)
		rootHandler = handler.NewMultiHandler(prettyHandler, jsonHandler)
	} else {
		rootHandler = jsonHandler
	}

	logger = slog.New(rootHandler)
}

func init() {
	opts = slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger = Default()
}

func isInteractive() bool {
	return term.IsTerminal(int(os.Stderr.Fd()))
}
