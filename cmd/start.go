package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/allan-simon/go-singleinstance"
	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/service"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the AIUB Notice Fetcher service",
	Long:  `Start the AIUB Notice Fetcher service to fetch and display notices from the AIUB website.`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			null, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			os.Stdout = null
			os.Stderr = null
		}

		checkInterval, err := cmd.Flags().GetDuration("interval")
		if err != nil {
			log.Fatalf("Error parsing interval flag: %v", err)
		}

		logFile, err := setupLogging()
		if err != nil {
			log.Fatalf("Error setting up logging: %v", err)
		}
		defer logFile.Close()

		lock, err := acquireLock()
		if err != nil {
			fmt.Println("Another instance is already running. Exiting.")
			return
		}
		defer lock.Close()

		log.Println("Single instance lock acquired.")
		service.Run(common.AUMID, checkInterval)
		log.Println("Service stopped. Exiting.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().DurationP("interval", "i", 30*time.Minute, "Set the interval for fetching notices")
	startCmd.Flags().Bool("quiet", false, "Suppress console output")
}

func setupLogging() (*os.File, error) {
	logPath, err := common.GetLogPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get log file path: %w", err)
	}

	logFile, err := os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	return logFile, nil
}

func acquireLock() (*os.File, error) {
	lockPath, _ := common.GetLockPath()
	return singleinstance.CreateLockFile(lockPath)
}
