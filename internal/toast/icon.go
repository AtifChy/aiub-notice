package toast

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AtifChy/aiub-notice/internal/common"
)

const iconURL = "https://www.aiub.edu/Files/Templates/AIUBv3/assets/images/aiub-logo-white-border.svg"

func fetchIcon(url, dest string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download icon: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download icon: received status code %d", response.StatusCode)
	}

	dir := filepath.Dir(dest)
	tmp, err := os.CreateTemp(dir, "aiub-icon-*.svg")
	if err != nil {
		return fmt.Errorf("failed to create temp icon file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, response.Body); err != nil {
		return fmt.Errorf("failed to write icon file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close temp icon file: %w", err)
	}

	if err := os.Rename(tmp.Name(), dest); err != nil {
		return fmt.Errorf("failed to move temp icon into place: %w", err)
	}

	return nil
}

func getIconPath() (string, error) {
	dataPath, err := common.GetDataPath()
	if err != nil {
		return "", fmt.Errorf("failed to get data path: %w", err)
	}

	iconPath := filepath.Join(dataPath, "aiub-icon.svg")
	iconPath, err = filepath.Abs(iconPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for icon: %w", err)
	}

	info, err := os.Stat(iconPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to stat icon file: %w", err)
	}

	if (info != nil && info.Size() <= 0) || os.IsNotExist(err) {
		if err = fetchIcon(iconURL, iconPath); err != nil {
			return "", fmt.Errorf("failed to download icon: %w", err)
		}
	}

	return iconPath, nil
}
