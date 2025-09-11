package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/appid"
	"github.com/AtifChy/aiub-notice/internal/common"
)

// appidCmd represents the aumid command
var appidCmd = &cobra.Command{
	Use:     "appid",
	Aliases: []string{"aumid"},
	Short:   "Manage appid registration for Windows notifications",
	Long: `The appid command allows you to register or unregister the application
AppUserModelID (appid) in the Windows registry. This is necessary for sending
toast notifications on Windows.

Examples:
	# Register the application appid
	aiub-notice appid --register

	# unregister the application appid
	aiub-notice appid --unregister`,
	Run: func(cmd *cobra.Command, args []string) {
		if register, _ := cmd.Flags().GetBool("register"); register {
			iconPath, err := common.GetIconPath()
			if err != nil {
				log.Fatalf("Error getting icon path: %v", err)
			}
			if err = appid.Register(common.AppID, common.DisplayName, iconPath); err != nil {
				log.Fatalf("Error registering toast application: %v", err)
			}
			log.Printf("Successfully registered AIUB Notice toast application with ID %s", common.AppID)
		} else if unregister, _ := cmd.Flags().GetBool("unregister"); unregister {
			appid.Unregister(common.AppID)
			log.Printf("Successfully unregistered AIUB Notice toast application with ID %s", common.AppID)
		} else {
			_ = cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(appidCmd)

	appidCmd.Flags().BoolP("register", "r", false, "Register application AUMID for toast notifications")
	appidCmd.Flags().BoolP("unregister", "d", false, "Unregister application AUMID for toast notifications")
}
