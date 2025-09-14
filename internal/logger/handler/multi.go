package handler

import (
	"context"
	"log/slog"
)

type MutltiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MutltiHandler {
	return &MutltiHandler{handlers: handlers}
}

func (h *MutltiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, hh := range h.handlers {
		if err := hh.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (h *MutltiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, hh := range h.handlers {
		if hh.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MutltiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		newHandlers[i] = hh.WithAttrs(attrs)
	}
	return &MutltiHandler{handlers: newHandlers}
}

func (h *MutltiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		newHandlers[i] = hh.WithGroup(name)
	}
	return &MutltiHandler{handlers: newHandlers}
}
