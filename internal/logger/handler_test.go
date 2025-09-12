package logger

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// captureStdout captures anything written to stdout during fn.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// normalizeOutput removes ANSI color codes and the dynamic leading timestamp.
func normalizeOutput(s string) string {
	s = strings.TrimSpace(s)
	s = ansiRE.ReplaceAllString(s, "")
	if i := strings.Index(s, "] "); i != -1 { // drop "[timestamp] "
		s = s[i+2:]
	}
	return s
}

func newHandlerNoSource() *PrettyHandler {
	return NewPrettyHandler(io.Discard, PrettyHandlerOptions{
		HandlerOptions: slog.HandlerOptions{AddSource: false},
	})
}

func TestPrettyHandler_InfoWithAttrs(t *testing.T) {
	h := newHandlerNoSource()

	r := slog.NewRecord(time.Date(2023, 5, 6, 7, 8, 9, 0, time.UTC), slog.LevelInfo, "hello world", 0)
	r.AddAttrs(slog.Int("a", 1), slog.String("b", "two words"))

	out := captureStdout(t, func() { _ = h.Handle(context.Background(), r) })
	got := normalizeOutput(out)
	want := `INFO: hello world a=1 b="two words"`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPrettyHandler_Levels(t *testing.T) {
	tests := []struct {
		lvl  slog.Level
		want string
	}{
		{slog.LevelDebug, "DEBUG: msg"},
		{slog.LevelInfo, "INFO: msg"},
		{slog.LevelWarn, "WARN: msg"},
		{slog.LevelError, "ERROR: msg"},
	}

	h := newHandlerNoSource()
	for _, tt := range tests {
		r := slog.NewRecord(time.Unix(0, 0), tt.lvl, "msg", 0)
		out := captureStdout(t, func() { _ = h.Handle(context.Background(), r) })
		got := normalizeOutput(out)
		if got != tt.want {
			t.Fatalf("level %v: got %q, want %q", tt.lvl, got, tt.want)
		}
	}
}

func TestPrettyHandler_NoAttrsNoExtraSpaces(t *testing.T) {
	h := newHandlerNoSource()
	r := slog.NewRecord(time.Unix(0, 0), slog.LevelInfo, "msg", 0)
	out := captureStdout(t, func() { _ = h.Handle(context.Background(), r) })
	got := normalizeOutput(out)
	want := "INFO: msg"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
