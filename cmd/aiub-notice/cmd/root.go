// Package cmd
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/logger"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     common.AppName,
	Short:   "AIUB Notice Notifier",
	Long:    `AIUB Notice Notifier is a command-line tool that fetches and displays notices from AIUB's official website.`,
	Version: common.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.L().Error("executing command", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
