// Package toast provides a simple way to show toast notifications for AIUB notices.
package toast

import (
	"fmt"

	"git.sr.ht/~jackmordaunt/go-toast/v2"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/notice"
)

func Show(notice notice.Notice) error {
	iconPath, err := getIconPath()
	if err != nil {
		return fmt.Errorf("failed to get icon path: %w", err)
	}

	notif := toast.Notification{
		AppID:               common.AppID,
		Title:               notice.Title,
		Body:                notice.Desc,
		Icon:                iconPath,
		ActivationType:      toast.Protocol,
		ActivationArguments: notice.Link,
		Actions: []toast.Action{
			{Type: toast.Protocol, Content: "Open", Arguments: notice.Link},
			{Type: toast.Protocol, Content: "Dismiss", Arguments: ""},
		},
	}

	return notif.Push()
}
