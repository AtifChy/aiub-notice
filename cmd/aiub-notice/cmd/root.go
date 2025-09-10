// Package cmd
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
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
		log.Fatalf("Error executing command: %v", err)
	}
}
