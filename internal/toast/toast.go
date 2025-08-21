// Package toast provides a simple way to show toast notifications for AIUB notices.
package toast

import (
	"github.com/electricbubble/go-toast"

	"github.com/AtifChy/aiub-notice/internal/notice"
)

func Show(aumid string, notice notice.Notice) error {
	return toast.Push(
		notice.Desc,
		toast.WithAppID(aumid),
		toast.WithTitle(notice.Title),
		toast.WithProtocolAction("Open", notice.Link),
		toast.WithProtocolAction("Dismiss"),
		toast.WithAudio(toast.Default),
	)
}
