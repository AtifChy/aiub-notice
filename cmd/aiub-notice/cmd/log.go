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
	RunE: func(cmd *cobra.Command, args []string) error {
		logPath := common.GetLogPath()

		if clear, _ := cmd.Flags().GetBool("clear"); clear {
			err := os.Truncate(logPath, 0)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf(`clearing log file: %w`, err)
			}
			logger.L().Info("log file cleared")
			return nil
		}

		logFile, err := os.Open(logPath)
		if os.IsNotExist(err) {
			logger.L().Info("log file does not exist", slog.String("path", logPath))
			return nil
		} else if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}
		defer logFile.Close()

		fmt.Printf("--- Displaying logs from %s ---\n", logPath)
		if _, err := io.Copy(os.Stdout, logFile); err != nil {
			return fmt.Errorf("reading log file: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().Bool("clear", false, "Clear the log file")
}
