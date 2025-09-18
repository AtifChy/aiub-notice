// Package cmd provides the command-line interface for the AIUB Notice Notifier application.
package cmd

import (
	"log/slog"

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
func Execute() error {
	logFile, err := common.GetLogFile()
	if err != nil {
		logger.L().Error("open log file", slog.String("error", err.Error()))
		logger.L().Warn("logging to file is disabled.")
	}
	if logFile != nil {
		logger.SetOutputFile(logFile)
		defer logFile.Close()
	}

	return rootCmd.Execute()
}
