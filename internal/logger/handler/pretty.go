package handler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path"
	"strconv"
	"strings"
	"time"

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
		dir := path.Base(path.Dir(source.File))
		file := path.Base(source.File)

		var sb strings.Builder
		sb.WriteString(dir)
		sb.WriteByte('/')
		sb.WriteString(file)
		sb.WriteByte(':')
		sb.WriteString(strconv.Itoa(source.Line))
		src = sb.String()
	}

	var attrsBuilder strings.Builder
	record.Attrs(func(a slog.Attr) bool {
		value := fmt.Sprintf("%v", a.Value)
		if strings.ContainsAny(value, " \t") {
			value = fmt.Sprintf("%q", value)
		}
		if attrsBuilder.Len() > 0 {
			attrsBuilder.WriteByte(' ')
		}
		attrsBuilder.WriteString(a.Key)
		attrsBuilder.WriteByte('=')
		attrsBuilder.WriteString(value)
		return true
	})

	var lineBuilder strings.Builder
	lineBuilder.WriteString(color.New(color.FgHiBlack).Sprint(ts))
	lineBuilder.WriteString(" [")
	lineBuilder.WriteString(levelColor.Sprint(record.Level.String()))
	lineBuilder.WriteString("] ")
	lineBuilder.WriteString(record.Message)
	if attrsBuilder.Len() > 0 {
		lineBuilder.WriteString(": ")
		lineBuilder.WriteString(attrsBuilder.String())
	}
	if src != "" {
		lineBuilder.WriteString(" (")
		lineBuilder.WriteString(color.New(color.FgHiBlack).Sprint(src))
		lineBuilder.WriteByte(')')
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
