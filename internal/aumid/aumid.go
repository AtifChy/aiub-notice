// Package aumid provides functions to register and unregister AppUserModelIDs (AUMIDs) in the Windows registry.
package aumid

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func Register(aumid, displayName, iconPath string) error {
	if iconPath != "" && !filepath.IsAbs(iconPath) {
		return fmt.Errorf("iconPath must be an absolute path: %s", iconPath)
	}

	regPath := fmt.Sprintf(`SOFTWARE\Classes\AppUserModelId\%s`, aumid)

	key, _, err := registry.CreateKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create registry key %s: %w", regPath, err)
	}
	defer key.Close()

	if err := key.SetStringValue("DisplayName", displayName); err != nil {
		return fmt.Errorf("failed to set DisplayName for %s: %w", aumid, err)
	}

	if iconPath != "" {
		if err := key.SetStringValue("IconUri", iconPath); err != nil {
			return fmt.Errorf("failed to set IconUri for %s: %w", aumid, err)
		}
	} else {
		// when iconPath is empty, remove any existing IconUri value
		_ = key.DeleteValue("IconUri")
	}

	return nil
}

func Unregister(aumid string) error {
	regPath := fmt.Sprintf(`SOFTWARE\Classes\AppUserModelId\%s`, aumid)

	err := registry.DeleteKey(registry.CURRENT_USER, regPath)
	if err != nil {
		return fmt.Errorf("failed to delete registry key %s: %w", regPath, err)
	}

	return nil
}
