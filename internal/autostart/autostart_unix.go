//go:build !windows

// Package autostart provides functionality to manage application autostart on Unix-like systems.
package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AtifChy/aiub-notice/internal/common"
)

func getStartupPath() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(config, "autostart"), nil
}

func EnableAutostart(interval time.Duration) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("get absolute path of executable: %w", err)
	}

	startupPath, err := getStartupPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(startupPath, 0o755); err != nil {
		return fmt.Errorf("create autostart directory: %w", err)
	}

	desktopEntry := fmt.Sprintf(
		`[Desktop Entry]
Type=Application
Name=AIUB Notice
Exec=%s start --interval %s
Version=%s
Hidden=false
X-GNOME-Autostart-enabled=true
Comment=Start AIUB Notice on login
`, exePath, interval, common.Version)

	desktopFilePath := filepath.Join(startupPath, "aiub-notice.desktop")

	return os.WriteFile(desktopFilePath, []byte(desktopEntry), 0o644)
}

func DisableAutostart() error {
	startupPath, err := getStartupPath()
	if err != nil {
		return err
	}

	desktopFilePath := filepath.Join(startupPath, "aiub-notice.desktop")
	err = os.Remove(desktopFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("remove autostart file: %w", err)
	}

	return nil
}

func IsAutostartEnabled() (bool, error) {
	startupPath, err := getStartupPath()
	if err != nil {
		return false, err
	}

	desktopFilePath := filepath.Join(startupPath, "aiub-notice.desktop")
	_, err = os.Stat(desktopFilePath)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
