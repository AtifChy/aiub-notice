package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/logger"
	"github.com/AtifChy/aiub-notice/internal/service"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"close"},
	Short:   "stop the AIUB Notice Fetcher service",
	Long:    `This command stops the AIUB Notice Fetcher service.`,
	Run: func(cmd *cobra.Command, args []string) {
		proc, err := service.GetProcessFromLock()
		if err != nil {
			logger.L().Error("retrieving service process", slog.String("error", err.Error()))
			return
		}
		if proc == nil {
			logger.L().Info("no running service found.")
			return
		}
		if err = proc.Signal(os.Kill); err != nil {
			logger.L().Error("sending interrupt signal", slog.String("error", err.Error()))
			return
		}
		logger.L().Info("service stopped successfully.")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
