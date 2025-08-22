package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/autostart"
)

// autostartCmd represents the autostart command
var autostartCmd = &cobra.Command{
	Use:   "autostart",
	Short: "Manage autostart settings for AIUB Notice Fetcher service",
	Long:  `This command allows you to enable or disable autostart for the AIUB Notice Fetcher service on Windows systems.`,
	Run: func(cmd *cobra.Command, args []string) {
		if enable, _ := cmd.Flags().GetBool("enable"); enable {
			interval, err := cmd.Flags().GetDuration("interval")
			if err != nil {
				log.Fatalf("Error parsing interval flag: %v", err)
			}

			err = autostart.EnableAutostart(interval)
			if err != nil {
				log.Fatalf("Error enabling autostart: %v", err)
			}

			fmt.Println("Autostart enabled for AIUB Notice Fetcher service.")
		} else if disable, _ := cmd.Flags().GetBool("disable"); disable {
			err := autostart.DisableAutostart()
			if err != nil {
				log.Fatalf("Error disabling autostart: %v", err)
			}

			fmt.Println("Autostart disabled for AIUB Notice Fetcher service.")
		} else if status, _ := cmd.Flags().GetBool("status"); status {
			enabled, err := autostart.IsAutostartEnabled()
			if err != nil {
				log.Fatalf("Error checking autostart status: %v", err)
			}

			fmt.Printf("Autostart is currently %s.\n", map[bool]string{true: "enabled", false: "disabled"}[enabled])
		} else {
			_ = cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(autostartCmd)

	autostartCmd.Flags().Bool("enable", false, "Enable autostart for the AIUB Notice Fetcher service")
	autostartCmd.Flags().DurationP("interval", "i", 30*time.Minute, "Set the interval for fetching notices [Used with --enable]")

	autostartCmd.Flags().Bool("disable", false, "Disable autostart for the AIUB Notice Fetcher service")
	autostartCmd.Flags().BoolP("status", "s", false, "Check if autostart is enabled")
}
