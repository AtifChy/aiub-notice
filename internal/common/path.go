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

var GetLogPath = func() string {
	return filepath.Join(os.TempDir(), AppName+".log")
}

// GetLogFile returns a file handle for the application's log file.
func GetLogFile() (*os.File, error) {
	const logSize = 5 * 1024 * 1024 // 5 MB
	logPath := GetLogPath()

	info, err := os.Stat(logPath)
	if err == nil && info.Size() > logSize {
		if err := os.Truncate(logPath, 0); err != nil {
			return nil, err
		}
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	logFile, err := os.OpenFile(
		logPath,
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
