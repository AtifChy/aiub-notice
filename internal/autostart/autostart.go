// Package autostart provides functions to manage the autostart behavior of the application on Windows.
package autostart

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AtifChy/aiub-notice/internal/common"
	"golang.org/x/sys/windows/registry"
)

const autostartKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`

func EnableAutostart() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of executable: %w", err)
	}

	autostartValue := fmt.Sprintf(`"%s" start`, exePath)

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		autostartKeyPath,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to create registry key: %w", err)
	}
	defer key.Close()

	return key.SetStringValue(common.AppName, autostartValue)
}

func DisableAutostart() error {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		autostartKeyPath,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	return key.DeleteValue(common.AppName)
}

func IsAutostartEnabled() (bool, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		autostartKeyPath,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false, fmt.Errorf("failed to open registry key: %w", err)
	}

	_, _, err = key.GetStringValue(common.AppName)
	if err == registry.ErrNotExist {
		return false, nil // Autostart is not enabled
	}
	if err != nil {
		return false, fmt.Errorf("failed to get autostart value: %w", err)
	}

	return true, nil // Autostart is enabled
}
