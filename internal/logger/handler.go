package logger

import (
	"fmt"
	"io"
	"log/slog"
	"path"
	"strings"
	"time"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	opts PrettyHandlerOptions
}

func (h *PrettyHandler) handle(record slog.Record) error {
	level := record.Level.String() + ":"

	switch record.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	var attrs []string
	record.Attrs(func(a slog.Attr) bool {
		value := fmt.Sprintf("%v", a.Value.Any())
		if strings.ContainsAny(value, " \t") {
			value = fmt.Sprintf("%q", value)
		}
		attrs = append(attrs, fmt.Sprintf("%s=%s", a.Key, value))
		return true
	})

	var sourceStr string
	if h.opts.AddSource && record.PC != 0 {
		source := record.Source()
		dir := path.Base(path.Dir(source.File))
		file := path.Base(source.File)
		sourceStr = fmt.Sprintf("%s/%s:%d", dir, file, source.Line)
	}

	timeStr := record.Time.Local().Format(time.DateTime)
	timeStr = color.HiBlackString("[%s]", timeStr)
	msg := color.New(color.Bold).Sprint(record.Message)

	parts := []string{timeStr, level}
	if sourceStr != "" {
		parts = append(parts, color.CyanString(sourceStr))
	}
	parts = append(parts, msg)
	if len(attrs) > 0 {
		s := strings.Join(attrs, " ")
		parts = append(parts, color.WhiteString(s))
	}

	fmt.Println(strings.Join(parts, " "))

	return nil
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.HandlerOptions),
		opts:    opts,
	}
}
