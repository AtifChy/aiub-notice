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

func GetLogPath() string {
	return filepath.Join(os.TempDir(), AppName+".log")
}

// GetLogFile returns a file handle for the application's log file.
func GetLogFile() (*os.File, error) {
	logFile, err := os.OpenFile(
		GetLogPath(),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		return nil, err
	}
	return logFile, nil
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
