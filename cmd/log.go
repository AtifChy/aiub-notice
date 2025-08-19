package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View the log of notices",
	Long:  `View the log of notices fetched by the AIUB Notice Fetcher.`,
	Run: func(cmd *cobra.Command, args []string) {
		logPath, err := common.GetLogPath()
		if err != nil {
			log.Fatalf("Error getting log file path: %v", err)
		}

		clear, err := cmd.Flags().GetBool("clear")
		if err != nil {
			log.Fatalf("Error getting clear flag: %v", err)
		}

		if clear {
			file, err := os.OpenFile(logPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Error clearing log file: %v", err)
			}
			file.Close()
			fmt.Println("Log file cleared.")
			return
		}

		logFile, err := os.Open(logPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No log file found. Please run the service first to generate logs.")
				os.Exit(0)
			}
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
