package main

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/logger"
)

func main() {
	status := run()
	os.Exit(status)
}

func run() int {
	logfile, _ := common.GetLogFile()
	defer logfile.Close()
	logger.SetOutputFile(logfile)

	exe, err := os.Executable()
	if err != nil {
		logger.L().Error("getting executable path", slog.String("error", err.Error()))
		return 1
	}

	dir := filepath.Dir(exe)
	realExe := filepath.Join(dir, common.AppName+".exe")

	cmd := exec.Command(realExe, os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Start(); err != nil {
		logger.L().Error("starting command", slog.String("error", err.Error()))
		return 1
	}

	return 0
}
