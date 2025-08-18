package common

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetDataPath returns the path to the application's data directory.
func GetDataPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user data directory: %w", err)
	}
	return filepath.Join(cacheDir, AppName), nil
}

// GetLogPath returns the path to the log file for the application.
func GetLogPath() (string, error) {
	logDir, err := GetDataPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}
	return filepath.Join(logDir, AppName+".log"), nil
}

// GetTempPath returns the path to the temporary directory for the application.
func GetTempPath() (string, error) {
	tempDir := filepath.Join(os.TempDir(), AppName)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	return tempDir, nil
}

// GetLockPath returns the path to the lock file for the application.
func GetLockPath() (string, error) {
	lockDir, err := GetTempPath()
	if err != nil {
		return "", fmt.Errorf("failed to get temporary directory: %w", err)
	}
	return filepath.Join(lockDir, AppName+".lock"), nil
}
