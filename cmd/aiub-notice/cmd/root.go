// Package cmd
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
)

var rootCmd = &cobra.Command{
	Use:     common.AppName,
	Short:   "AIUB Notice Notifier",
	Long:    `AIUB Notice Notifier is a command-line tool that fetches and displays notices from AIUB's official website.`,
	Version: common.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
