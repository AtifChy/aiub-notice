// Package register provides functionality to register an application with the Windows registry
package register

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func Register(aumid, displayName, iconPath string) error {
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
		if !filepath.IsAbs(iconPath) {
			return fmt.Errorf("iconPath must be an absolute path: %s", iconPath)
		}
		if err := key.SetStringValue("IconUri", iconPath); err != nil {
			return fmt.Errorf("failed to set IconUri for %s: %w", aumid, err)
		}
	} else {
		_ = key.DeleteValue("IconUri")
	}

	return nil
}
