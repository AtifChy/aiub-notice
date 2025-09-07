package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/service"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"info"},
	Short:   "Check the status of the AIUB Notice Fetcher service",
	Long:    `This command checks whether the AIUB Notice Fetcher service is currently running or not.`,
	Run: func(cmd *cobra.Command, args []string) {
		var isRunning bool
		proc, _ := service.GetProcessFromLock()
		if proc != nil {
			isRunning = true
		}
		fmt.Printf("AIUB Notice Fetcher service is currently %s.\n", map[bool]string{true: "running", false: "not running"}[isRunning])
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
