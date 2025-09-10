package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View the log of notices",
	Long:  `View the log of notices fetched by the AIUB Notice Fetcher.`,
	Run: func(cmd *cobra.Command, args []string) {
		logPath := common.GetLogPath()

		if clear, _ := cmd.Flags().GetBool("clear"); clear {
			err := os.Truncate(logPath, 0)
			if err != nil && !os.IsNotExist(err) {
				log.Fatalf("Error clearing log file: %v", err)
			} else {
				fmt.Println("Log file cleared.")
			}
			return
		}

		logFile, err := os.Open(logPath)
		if os.IsNotExist(err) {
			fmt.Println("No log file found. Please run the service first to generate logs.")
			os.Exit(0)
		} else if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer logFile.Close()

		fmt.Printf("--- Displaying logs from %s ---\n", logPath)
		if _, err := io.Copy(os.Stdout, logFile); err != nil {
			log.Fatalf("Error reading log file: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().Bool("clear", false, "Clear the log file")
}
