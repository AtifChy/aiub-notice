package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/logger"
	"github.com/AtifChy/aiub-notice/internal/logger/handler"
)

var showSource bool

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View the log of notices",
	Long:  `View the log of notices fetched by the AIUB Notice Fetcher.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logPath := common.GetLogPath()

		if clear, _ := cmd.Flags().GetBool("clear"); clear {
			err := os.Truncate(logPath, 0)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf(`clearing log file: %w`, err)
			}
			logger.L().Info("log file cleared")
			return nil
		}

		if source, _ := cmd.Flags().GetBool("source"); source {
			showSource = true
		}

		logFile, err := os.Open(logPath)
		if os.IsNotExist(err) {
			logger.L().Info("log file does not exist", slog.String("path", logPath))
			return nil
		} else if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}
		defer func() { _ = logFile.Close() }()

		fmt.Printf("--- Displaying logs from %s ---\n", logPath)

		return replayLogs(logFile)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().Bool("clear", false, "Clear the log file")
	logCmd.Flags().Bool("source", false, "Show source information in logs")
}

func replayLogs(logFile io.Reader) error {
	var err error
	out := slog.New(handler.NewPrettyHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	))
	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {
		var obj map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &obj); err != nil {
			continue
		}

		var datetime time.Time
		if t, ok := obj[slog.TimeKey].(string); ok {
			datetime, err = time.Parse(time.RFC3339, t)
			if err != nil {
				return fmt.Errorf("parsing time from log: %w", err)
			}
		}

		var level slog.Level
		if l, ok := obj[slog.LevelKey].(slog.Level); ok {
			level = l
		}

		var msg string
		if m, ok := obj[slog.MessageKey].(string); ok {
			msg = m
		}

		r := slog.Record{
			Time:    datetime,
			Level:   level,
			Message: msg,
		}

		for k, v := range obj {
			if k == slog.TimeKey || k == slog.LevelKey || k == slog.MessageKey {
				continue
			} else if k == slog.SourceKey && !showSource {
				continue
			}
			r.AddAttrs(slog.Any(k, v))
		}

		err = out.Handler().Handle(context.Background(), r)
		if err != nil {
			return fmt.Errorf("handling log record: %w", err)
		}
	}

	return scanner.Err()
}
