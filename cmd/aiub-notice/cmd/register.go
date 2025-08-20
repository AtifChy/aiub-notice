package cmd

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/register"
	"github.com/AtifChy/aiub-notice/internal/toast"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register AIUB Notice Fetcher with Windows Toast Notifications (Recommended)",
	Long:  `This command registers the AIUB Notice Fetcher application with Windows Toast Notifications, allowing it to display notifications in the system tray.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if err = register.Register(common.AUMID, "AIUB Notice", iconPath); err != nil {
			log.Fatalf("Error registering toast application: %v", err)
		}

		log.Printf("Successfully registered AIUB Notice toast application with ID %s", common.AUMID)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
