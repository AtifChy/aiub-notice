package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/AtifChy/aiub-notice/internal/autostart"
	"github.com/spf13/cobra"
)

// autostartCmd represents the autostart command
var autostartCmd = &cobra.Command{
	Use:   "autostart",
	Short: "Manage autostart settings for AIUB Notice Fetcher service",
	Long:  `This command allows you to enable or disable autostart for the AIUB Notice Fetcher service on Windows systems.`,
	Run: func(cmd *cobra.Command, args []string) {
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")
		status, _ := cmd.Flags().GetBool("status")

		if enable {
			interval, err := cmd.Flags().GetDuration("interval")
			if err != nil {
				log.Fatalf("Error parsing interval flag: %v", err)
			}

			err = autostart.EnableAutostart(interval)
			if err != nil {
				log.Fatalf("Error enabling autostart: %v", err)
			}

			fmt.Println("Autostart enabled for AIUB Notice Fetcher service.")
		} else if disable {
			err := autostart.DisableAutostart()
			if err != nil {
				log.Fatalf("Error disabling autostart: %v", err)
			}
			fmt.Println("Autostart disabled for AIUB Notice Fetcher service.")
		} else if status {
			enabled, err := autostart.IsAutostartEnabled()
			if err != nil {
				log.Fatalf("Error checking autostart status: %v", err)
			}
			fmt.Printf("Autostart is currently %s.\n", map[bool]string{true: "enabled", false: "disabled"}[enabled])
		} else {
			fmt.Println("Please specify either --enable or --disable flag.")
		}
	},
}

func init() {
	rootCmd.AddCommand(autostartCmd)

	autostartCmd.Flags().Bool("enable", false, "Enable autostart for the AIUB Notice Fetcher service")
	autostartCmd.Flags().DurationP("interval", "i", 1*time.Hour, "Set the interval for fetching notices [Used with --enable]")

	autostartCmd.Flags().Bool("disable", false, "Disable autostart for the AIUB Notice Fetcher service")
	autostartCmd.Flags().BoolP("status", "s", false, "Check if autostart is enabled")
}
