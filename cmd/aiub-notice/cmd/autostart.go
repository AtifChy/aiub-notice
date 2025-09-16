package cmd

import (
	"fmt"
	"log/slog"
	"os"
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
	Run: func(cmd *cobra.Command, args []string) {
		if enable, _ := cmd.Flags().GetBool("enable"); enable {
			interval, err := cmd.Flags().GetDuration("interval")
			if err != nil {
				logger.L().Error("parsing interval flag", slog.String("error", err.Error()))
				os.Exit(1)
			}

			err = autostart.EnableAutostart(interval)
			if err != nil {
				logger.L().Error("enabling autostart", slog.String("error", err.Error()))
				os.Exit(1)
			}

			logger.L().Info("Autostart enabled for AIUB Notice Fetcher service.")
		} else if disable, _ := cmd.Flags().GetBool("disable"); disable {
			err := autostart.DisableAutostart()
			if err != nil {
				logger.L().Error("disabling autostart", slog.String("error", err.Error()))
				os.Exit(1)
			}

			logger.L().Info("Autostart disabled for AIUB Notice Fetcher service.")
		} else if status, _ := cmd.Flags().GetBool("status"); status {
			enabled, err := autostart.IsAutostartEnabled()
			if err != nil {
				logger.L().Error("checking autostart status", slog.String("error", err.Error()))
				os.Exit(1)
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
