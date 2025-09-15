package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/allan-simon/go-singleinstance"
	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/logger"
	"github.com/AtifChy/aiub-notice/internal/service"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "Start the AIUB Notice Fetcher service",
	Long:    `Start the AIUB Notice Fetcher service to fetch and display notices from the AIUB website.`,
	Run:     run,
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().DurationP("interval", "i", 30*time.Minute, "Set the interval for fetching notices")
	startCmd.Flags().Bool("quiet", false, "Suppress console output")
}

func run(cmd *cobra.Command, args []string) {
	quiet, _ := cmd.Flags().GetBool("quiet")
	if quiet {
		null, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		os.Stdout = null
		os.Stderr = null
	}

	logfile, err := common.GetLogFile()
	if err != nil {
		logger.L().Error("open log file", slog.String("error", err.Error()))
		logger.L().Warn("logging to file is disabled.")
	}
	if logfile != nil {
		logger.SetOutputFile(logfile)
		defer logfile.Close()
	}

	lock, err := acquireLock()
	if err != nil {
		logger.L().Info("another instance is already running, exiting...")
		return
	}
	logger.L().Info("single instance lock acquired.")
	defer lock.Close()

	checkInterval, err := cmd.Flags().GetDuration("interval")
	if err != nil {
		logger.L().Error("parsing interval flag", slog.String("error", err.Error()))
		return
	}
	service.Run(checkInterval)

	logger.L().Info("service stopped.")
}

func acquireLock() (*os.File, error) {
	lockPath, _ := common.GetLockPath()
	return singleinstance.CreateLockFile(lockPath)
}
