package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/appid"
	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/logger"
)

// appidCmd represents the appid command
var appidCmd = &cobra.Command{
	Use:     "appid",
	Aliases: []string{"aumid"},
	Short:   "Manage AppID registration for Windows notifications",
	Long: `The appid command allows you to register or unregister the application
AppUserModelID (AUMID) in the Windows registry. This is necessary for sending
toast notifications on Windows.

Examples:
	# Register the application appid
	aiub-notice appid --register

	# unregister the application appid
	aiub-notice appid --unregister`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if register, _ := cmd.Flags().GetBool("register"); register {
			iconPath, err := common.GetIconPath()
			if err != nil {
				return fmt.Errorf("getting icon path: %w", err)
			}
			if err = appid.Register(common.AppID, common.DisplayName, iconPath); err != nil {
				return fmt.Errorf("registering appid: %w", err)
			}
			logger.L().Info("successfully registered AIUB Notice toast application", slog.String("appid", common.AppID))
		} else if unregister, _ := cmd.Flags().GetBool("unregister"); unregister {
			appid.Unregister(common.AppID)
			logger.L().Info("successfully unregistered AIUB Notice toast application", slog.String("appid", common.AppID))
		} else {
			_ = cmd.Help()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(appidCmd)

	appidCmd.Flags().BoolP("register", "r", false, "Register application AppID for toast notifications")
	appidCmd.Flags().BoolP("unregister", "d", false, "Unregister application AppID for toast notifications")
}
