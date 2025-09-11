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

func fetchIcon(url, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download icon: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download icon: received status code %d", response.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create icon file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("failed to write icon file: %w", err)
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

	if _, err := os.Stat(iconPath); err != nil {
		if os.IsNotExist(err) {
			err = fetchIcon(iconURL, iconPath)
			if err != nil {
				return "", fmt.Errorf("failed to download icon: %w", err)
			}
		}
	}

	return iconPath, nil
}
