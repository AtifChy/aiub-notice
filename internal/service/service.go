// Package service provides the main application logic for checking and notifying about new notices.
package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/logger"
	"github.com/AtifChy/aiub-notice/internal/notice"
	"github.com/AtifChy/aiub-notice/internal/toast"
)

// Run starts the notice checking service.
func Run(checkInterval time.Duration) {
	logger.L().Info("starting initial notice check...")

	// Load previously seen notices
	seenNotices, err := notice.LoadSeenNotices()
	if err != nil {
		logger.L().Error("loading seen notices", slog.String("error", err.Error()))
		seenNotices = make(map[string]struct{})
	}

	// Perform initial check for notices
	if err = checkNotice(seenNotices); err != nil {
		logger.L().Error(
			"initial notice check",
			slog.String("error", err.Error()),
		)
	}

	// Context for graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	// Start ticker for periodic checks
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	logger.L().Info("service started", slog.String("check_interval", checkInterval.String()))

	// Main service loop
	for {
		select {
		case <-ticker.C:
			logger.L().Info("checking for new notices...")
			if err := checkNotice(seenNotices); err != nil {
				logger.L().Error("checking for new notices", slog.String("error", err.Error()))
			}

		case <-ctx.Done():
			logger.L().Info("received shutdown signal, stopping service...")
			return
		}
	}
}

func checkNotice(seenNotices map[string]struct{}) error {
	notices, err := notice.GetNotices()
	if err != nil {
		return fmt.Errorf("fetch notices: %w", err)
	}

	var newNotices []notice.Notice
	for _, n := range notices {
		if _, seen := seenNotices[n.Link]; !seen {
			newNotices = append(newNotices, n)
			seenNotices[n.Link] = struct{}{}
		}
	}

	if len(newNotices) > 0 {
		logger.L().Info("found new notices", slog.Int("count", len(newNotices)))

		path, err := notice.GetSeenNoticesPath()
		if err != nil {
			logger.L().Error("getting seen notices path", slog.String("error", err.Error()))
		}

		if _, err = os.Stat(path); err == nil {
			for _, n := range newNotices {
				err := toast.Show(n)
				if err != nil {
					logger.L().Error(
						"showing toast notification",
						slog.String("title", n.Title),
						slog.String("error", err.Error()),
					)
				} else {
					logger.L().Info("sent notification for notice", slog.String("title", n.Title))
				}
			}
		} else if os.IsNotExist(err) {
			logger.L().Warn("seen notices file does not exist, skipping notifications")
		} else {
			logger.L().Error("checking seen notices file", slog.String("error", err.Error()))
		}

		if err := notice.SaveSeenNotices(seenNotices); err != nil {
			return fmt.Errorf("save seen notices: %w", err)
		}
	} else {
		logger.L().Info("no new notices found.")
	}

	return nil
}

func GetProcessFromLock() (*os.Process, error) {
	lockPath, err := common.GetLockPath()
	if err != nil {
		return nil, fmt.Errorf("get lock file path: %w", err)
	}

	data, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("read lock file: %w", err)
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PID in lock file: %w", err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("find process with PID %d: %w", pid, err)
	}

	return proc, nil
}
