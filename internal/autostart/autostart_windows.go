//go:build windows

// Package autostart provides functions to manage the autostart behavior of the application on Windows.
package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jxeng/shortcut"

	"github.com/AtifChy/aiub-notice/internal/common"
)

func getStartupPath() (string, error) {
	appData, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get appdata directory: %w", err)
	}
	startupPath := filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	return startupPath, nil
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
	exeDir := filepath.Dir(exePath)

	launcherPath := filepath.Join(exeDir, common.LauncherName+".exe")
	args := fmt.Sprintf("start --interval %s", interval)

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
		return fmt.Errorf("remove autostart file: %w", err)
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
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
