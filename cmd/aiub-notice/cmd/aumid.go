package cmd

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/aumid"
	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/toast"
)

// aumidCmd represents the aumid command
var aumidCmd = &cobra.Command{
	Use:     "aumid",
	Aliases: []string{"appid"},
	Short:   "Manage AUMID registration for Windows notifications",
	Long: `The aumid command allows you to register or deregister the application
AppUserModelID (AUMID) in the Windows registry. This is necessary for sending
toast notifications on Windows.

Examples:
	# Register the application AUMID
	aiub-notice aumid --register

	# Deregister the application AUMID
	aiub-notice aumid --deregister`,
	Run: func(cmd *cobra.Command, args []string) {
		if register, _ := cmd.Flags().GetBool("register"); register {
			iconURL := "https://www.aiub.edu/Files/Templates/AIUBv3/assets/images/aiub-logo-white-border.svg"

			dataPath, err := common.GetDataPath()
			if err != nil {
				log.Fatalf("Error getting data path: %v", err)
			}

			iconPath := filepath.Join(dataPath, "aiub-icon.svg")
			iconPath, err = filepath.Abs(iconPath)
			if err != nil {
				log.Fatalf("Error getting absolute path for icon: %v", err)
			}

			if err = toast.DownloadIcon(iconURL, iconPath); err != nil {
				log.Fatalf("Error downloading icon: %v", err)
			}

			if err = aumid.Register(common.AUMID, "AIUB Notice", iconPath); err != nil {
				log.Fatalf("Error registering toast application: %v", err)
			}

			log.Printf("Successfully registered AIUB Notice toast application with ID %s", common.AUMID)
		} else if deregister, _ := cmd.Flags().GetBool("deregister"); deregister {
			aumid.Deregister(common.AUMID)
			log.Printf("Successfully deregistered AIUB Notice toast application with ID %s", common.AUMID)
		} else {
			_ = cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(aumidCmd)

	aumidCmd.Flags().BoolP("register", "r", false, "Register application AUMID for toast notifications")
	aumidCmd.Flags().BoolP("deregister", "d", false, "Deregister application AUMID for toast notifications")
}
