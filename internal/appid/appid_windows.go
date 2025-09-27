// Package appid provides functions to register and unregister AppUserModelIDs (AUMIDs) in the Windows registry.
package appid

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const appKeyRoot = `SOFTWARE\Classes\AppUserModelId`

func Register(appID, displayName, iconPath string) error {
	if iconPath != "" && !filepath.IsAbs(iconPath) {
		return fmt.Errorf("iconPath must be an absolute path: %s", iconPath)
	}

	regPath := appKeyRoot + `\` + appID

	key, _, err := registry.CreateKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("create registry key %s: %w", regPath, err)
	}
	defer func() { _ = key.Close() }()

	if err := key.SetStringValue("DisplayName", displayName); err != nil {
		return fmt.Errorf("set DisplayName for %s: %w", appID, err)
	}

	if iconPath != "" {
		if err := key.SetStringValue("IconUri", iconPath); err != nil {
			return fmt.Errorf("set IconUri for %s: %w", appID, err)
		}
	} else {
		// when iconPath is empty, remove any existing IconUri value
		_ = key.DeleteValue("IconUri")
	}

	return nil
}

func Unregister(appID string) error {
	regPath := appKeyRoot + `\` + appID

	err := registry.DeleteKey(registry.CURRENT_USER, regPath)
	if err != nil {
		return fmt.Errorf("delete registry key %s: %w", regPath, err)
	}

	return nil
}
