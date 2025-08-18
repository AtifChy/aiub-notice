package toast

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadIcon(url, path string) error {
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
