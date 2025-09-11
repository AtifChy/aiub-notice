package notice

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AtifChy/aiub-notice/internal/common"
)

func getCachedNoticesPath() (string, error) {
	path, err := common.GetDataPath()
	if err != nil {
		return "", fmt.Errorf("get data path: %w", err)
	}
	return filepath.Join(path, "notices.json"), nil
}

func storeCachedNotices(notices []Notice) error {
	path, err := getCachedNoticesPath()
	if err != nil {
		return fmt.Errorf("get cache path: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create cache file: %w", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(notices)
}

func GetCachedNotices() ([]Notice, error) {
	path, err := getCachedNoticesPath()
	if err != nil {
		return nil, fmt.Errorf("get cache path: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open cache file: %w", err)
	}
	defer file.Close()

	var notices []Notice
	if err := json.NewDecoder(file).Decode(&notices); err != nil {
		return nil, fmt.Errorf("decode cache file: %w", err)
	}

	return notices, nil
}

// GetSeenNoticesPath returns the path to the file where seen notices are stored.
func GetSeenNoticesPath() (string, error) {
	path, err := common.GetDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "seen_notices.json"), nil
}

func LoadSeenNotices() (map[string]struct{}, error) {
	path, err := GetSeenNoticesPath()
	if err != nil {
		return nil, fmt.Errorf("get seen notices file path: %w", err)
	}
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return make(map[string]struct{}), nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var seen map[string]struct{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&seen); err != nil {
		return make(map[string]struct{}), nil // fallback to empty map if decoding fails
	}
	return seen, nil
}

func SaveSeenNotices(seen map[string]struct{}) error {
	path, err := GetSeenNoticesPath()
	if err != nil {
		return fmt.Errorf("get seen notices file path: %w", err)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(seen)
}
