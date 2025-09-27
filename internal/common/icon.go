package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const iconURL = "https://www.aiub.edu/Files/Templates/AIUBv3/assets/images/aiub-logo-white-border.svg"

func fetchIcon(url, dest string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download icon: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download icon: received status code %d", response.StatusCode)
	}

	dir := filepath.Dir(dest)
	tmp, err := os.CreateTemp(dir, "aiub-icon-*.svg")
	if err != nil {
		return fmt.Errorf("create temp icon file: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()

	if _, err := io.Copy(tmp, response.Body); err != nil {
		return fmt.Errorf("write icon file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp icon file: %w", err)
	}

	if err := os.Rename(tmp.Name(), dest); err != nil {
		return fmt.Errorf("move temp icon into place: %w", err)
	}

	return nil
}

// ensureIconExists checks if the icon file exists at the given path,
// and downloads it if it does not exist or is empty.
func ensureIconExists(path string) (string, error) {
	iconPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("get icon absolute path: %w", err)
	}

	info, err := os.Stat(iconPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("stat icon file: %w", err)
	}

	if (info != nil && info.Size() <= 0) || os.IsNotExist(err) {
		if err = fetchIcon(iconURL, iconPath); err != nil {
			return "", fmt.Errorf("download icon: %w", err)
		}
	}

	return iconPath, nil
}
