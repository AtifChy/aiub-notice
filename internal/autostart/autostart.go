// Package autostart provides functions to manage the autostart behavior of the application on Windows.
package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/jxeng/shortcut"
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
	exeDir := filepath.Dir(exePath)

	launcherPath := filepath.Join(exeDir, common.LauncherName+".exe")
	args := fmt.Sprintf("start --interval %s --quiet", interval)

	startupPath, err := getStartupPath()
	if err != nil {
		return err
	}

	sc := shortcut.Shortcut{
		ShortcutPath: filepath.Join(startupPath, common.AppName+".lnk"),
		Target:       launcherPath,
		IconLocation: launcherPath + ",0",
		Arguments:    args,
	}

	return shortcut.Create(sc)
}

func DisableAutostart() error {
	startupPath, err := getStartupPath()
	if err != nil {
		return err
	}

	scPath := filepath.Join(startupPath, common.AppName+".lnk")
	if err = os.Remove(scPath); os.IsExist(err) {
		return fmt.Errorf("failed to remove autostart file: %w", err)
	}

	return nil
}

func IsAutostartEnabled() (bool, error) {
	startupPath, err := getStartupPath()
	if err != nil {
		return false, err
	}

	scPath := filepath.Join(startupPath, common.AppName+".lnk")
	_, err = os.Stat(scPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
