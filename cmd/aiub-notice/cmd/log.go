package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/logger"
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
				logger.L().Error("clearing log file", slog.String("error", err.Error()))
				os.Exit(1)
			} else {
				logger.L().Info("log file cleared")
			}
			return
		}

		logFile, err := os.Open(logPath)
		if os.IsNotExist(err) {
			fmt.Println("No log file found. Please run the service first to generate logs.")
			os.Exit(0)
		} else if err != nil {
			logger.L().Error("opening log file", slog.String("error", err.Error()))
		}
		defer logFile.Close()

		fmt.Printf("--- Displaying logs from %s ---\n", logPath)
		if _, err := io.Copy(os.Stdout, logFile); err != nil {
			logger.L().Error("reading log file", slog.String("error", err.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().Bool("clear", false, "Clear the log file")
}
