// Package autostart provides functions to manage the autostart behavior of the application on Windows.
package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AtifChy/aiub-notice/internal/common"
)

func getStartupPath() (string, error) {
	appData, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get appdata directory: %w", err)
	}
	startupPath := filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	return startupPath, nil
}

func EnableAutostart(interval time.Duration) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of executable: %w", err)
	}

	startupPath, err := getStartupPath()
	if err != nil {
		return err
	}

	batPath := filepath.Join(startupPath, common.AppName+".bat")
	batContent := fmt.Sprintf("@echo off\nstart /b \"\" \"%s\" start --interval %s --quiet\n", exePath, interval)
	return os.WriteFile(batPath, []byte(batContent), 0644)
}

func DisableAutostart() error {
	startupPath, err := getStartupPath()
	if err != nil {
		return err
	}

	batPath := filepath.Join(startupPath, common.AppName+".bat")
	if err = os.Remove(batPath); os.IsExist(err) {
		return fmt.Errorf("failed to remove autostart file: %w", err)
	}

	return nil
}

func IsAutostartEnabled() (bool, error) {
	startupPath, err := getStartupPath()
	if err != nil {
		return false, err
	}

	batPath := filepath.Join(startupPath, common.AppName+".bat")
	_, err = os.Stat(batPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
