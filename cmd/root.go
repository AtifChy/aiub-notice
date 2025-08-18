// Package cmd
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aiub-notice",
	Short: "AIUB Notice Notifier",
	Long:  `AIUB Notice Notifier is a command-line tool that fetches and displays notices from AIUB's official website.`,
	// Run:   runService,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
