// Package toast provides a simple way to show toast notifications for AIUB notices.
package toast

import (
	"git.sr.ht/~jackmordaunt/go-toast/v2"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/notice"
)

func Show(notice notice.Notice) error {
	notif := toast.Notification{
		AppID:               common.AppID,
		Title:               notice.Title,
		Body:                notice.Desc,
		ActivationType:      toast.Protocol,
		ActivationArguments: notice.Link,
		Actions: []toast.Action{
			{Type: toast.Protocol, Content: "Open", Arguments: notice.Link},
			{Type: toast.Protocol, Content: "Dismiss", Arguments: ""},
		},
	}

	return notif.Push()
}
