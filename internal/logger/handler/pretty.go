package handler

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/fatih/color"
)

type PrettyHandler struct {
	opts slog.HandlerOptions
	w    io.Writer
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		opts: *opts,
		w:    w,
	}
}

func (h *PrettyHandler) Handle(_ context.Context, record slog.Record) error {
	levelColor := color.New(color.FgWhite)

	switch record.Level {
	case slog.LevelDebug:
		levelColor = color.New(color.FgMagenta)
	case slog.LevelInfo:
		levelColor = color.New(color.FgBlue)
	case slog.LevelWarn:
		levelColor = color.New(color.FgYellow)
	case slog.LevelError:
		levelColor = color.New(color.FgRed)
	}

	ts := record.Time.Format(time.DateTime)

	var src string
	if h.opts.AddSource && record.PC != 0 {
		source := record.Source()

		idx := strings.Index(source.File, common.AppName)
		file := source.File
		if idx != -1 {
			file = source.File[idx+len(common.AppName)+1:]
		}

		var sb strings.Builder
		sb.WriteString(file)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(source.Line))
		src = sb.String()
	}

	attrs := make(map[string]any, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	attrsJSON, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return err
	}

	var lineBuilder strings.Builder

	// timestamp
	lineBuilder.WriteString(color.HiBlackString(ts))

	// level
	lineBuilder.WriteByte(' ')
	levelColor.Fprintf(&lineBuilder, "[%s]", record.Level.String())

	// message
	lineBuilder.WriteByte(' ')
	lineBuilder.WriteString(record.Message)

	// source
	if src != "" {
		lineBuilder.WriteByte(' ')
		color.New(color.FgHiBlack).Fprintf(&lineBuilder, "@%s", src)
	}

	// attributes
	if len(attrs) > 0 {
		lineBuilder.WriteByte(' ')
		lineBuilder.WriteString(string(attrsJSON))
	}

	lineBuilder.WriteByte('\n')

	_, _ = io.WriteString(h.w, lineBuilder.String())

	return nil
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return h
}
