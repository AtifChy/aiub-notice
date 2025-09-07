package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/service"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"close"},
	Short:   "stop the AIUB Notice Fetcher service",
	Long:    `This command stops the AIUB Notice Fetcher service.`,
	Run: func(cmd *cobra.Command, args []string) {
		proc, err := service.GetProcessFromLock()
		if err != nil {
			fmt.Println("Error retrieving service process:", err)
			return
		}
		if proc == nil {
			fmt.Println("No running service found.")
			return
		}
		if err = proc.Signal(os.Kill); err != nil {
			fmt.Println("Error sending interrupt signal:", err)
			return
		}
		fmt.Println("Service stopped successfully.")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
