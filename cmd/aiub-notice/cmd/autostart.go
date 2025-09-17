package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/autostart"
	"github.com/AtifChy/aiub-notice/internal/logger"
)

// autostartCmd represents the autostart command
var autostartCmd = &cobra.Command{
	Use:     "autostart",
	Aliases: []string{"startup"},
	Short:   "Manage autostart settings for AIUB Notice Fetcher service",
	Long:    `This command allows you to enable or disable autostart for the AIUB Notice Fetcher service on Windows systems.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if enable, _ := cmd.Flags().GetBool("enable"); enable {
			interval, err := cmd.Flags().GetDuration("interval")
			if err != nil {
				return fmt.Errorf("parsing interval flag: %w", err)
			}

			err = autostart.EnableAutostart(interval)
			if err != nil {
				return fmt.Errorf("enabling autostart: %w", err)
			}

			logger.L().Info("Autostart enabled for AIUB Notice Fetcher service.")
		} else if disable, _ := cmd.Flags().GetBool("disable"); disable {
			err := autostart.DisableAutostart()
			if err != nil {
				return fmt.Errorf("disabling autostart: %w", err)
			}

			logger.L().Info("Autostart disabled for AIUB Notice Fetcher service.")
		} else if status, _ := cmd.Flags().GetBool("status"); status {
			enabled, err := autostart.IsAutostartEnabled()
			if err != nil {
				return fmt.Errorf("checking autostart status: %w", err)
			}

			fmt.Printf("Autostart is currently %s.\n", map[bool]string{true: "enabled", false: "disabled"}[enabled])
		}

		_ = cmd.Help()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(autostartCmd)

	autostartCmd.Flags().Bool("enable", false, "Enable autostart for the AIUB Notice Fetcher service")
	autostartCmd.Flags().DurationP("interval", "i", 30*time.Minute, "Set the interval for fetching notices [Used with --enable]")

	autostartCmd.Flags().Bool("disable", false, "Disable autostart for the AIUB Notice Fetcher service")
	autostartCmd.Flags().BoolP("status", "s", false, "Check if autostart is enabled")
}
