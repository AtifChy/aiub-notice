// Package service provides the main application logic for checking and notifying about new notices.
package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/notice"
	"github.com/AtifChy/aiub-notice/internal/toast"
)

// Run starts the notice checking service.
func Run(aumid string, checkInterval time.Duration) {
	log.Println("Starting initial notice check...")

	// Load previously seen notices
	seenNotices, err := notice.LoadSeenNotices()
	if err != nil {
		log.Printf("Error loading seen notices: %v", err)
		seenNotices = make(map[string]struct{})
	}

	// Perform initial check for notices
	if err = checkNotice(aumid, seenNotices); err != nil {
		log.Printf("Error during initial notice check: %v", err)
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

	log.Printf("Starting notice check loop with interval: %s", checkInterval)

	// Main service loop
	for {
		select {
		case <-ticker.C:
			log.Println("Checking for new notices...")
			if err := checkNotice(aumid, seenNotices); err != nil {
				log.Printf("Error checking for new notices: %v", err)
			}

		case <-ctx.Done():
			log.Println("Received shutdown signal, stopping service...")
			return
		}
	}
}

func checkNotice(aumid string, seenNotices map[string]struct{}) error {
	notices, err := notice.GetNotices()
	if err != nil {
		return fmt.Errorf("failed to fetch notices: %w", err)
	}

	var newNotices []notice.Notice
	for _, n := range notices {
		if _, seen := seenNotices[n.Link]; !seen {
			newNotices = append(newNotices, n)
			seenNotices[n.Link] = struct{}{}
		}
	}

	if len(newNotices) > 0 {
		log.Printf("Found %d new notices. Sending notifications...", len(newNotices))

		path, err := notice.GetSeenNoticesPath()
		if err != nil {
			log.Printf("Error getting seen notices path: %v", err)
		}

		if _, err = os.Stat(path); err == nil {
			for _, n := range newNotices {
				err := toast.Show(aumid, n)
				if err != nil {
					log.Printf("Error showing toast notification for notice %s: %v", n.Title, err)
				}
				log.Printf("Sent notification for notice: '%s'", n.Title)
			}
		} else if os.IsNotExist(err) {
			log.Println("Seen notices file does not exist, skipping notifications.")
		} else {
			log.Printf("Error checking seen notices file: %v", err)
		}

		if err := notice.SaveSeenNotices(seenNotices); err != nil {
			return fmt.Errorf("failed to save seen notices: %w", err)
		}
	} else {
		log.Println("No new notices found.")
	}

	return nil
}

func GetProcessFromLock() (*os.Process, error) {
	lockPath, err := common.GetLockPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get lock file path: %w", err)
	}

	data, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PID in lock file: %w", err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("failed to find process with PID %d: %w", pid, err)
	}

	return proc, nil
}
