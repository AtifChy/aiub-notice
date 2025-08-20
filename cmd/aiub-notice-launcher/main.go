package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/AtifChy/aiub-notice/internal/common"
)

func main() {
	logfile, _ := common.GetLogFile()
	defer logfile.Close()
	log.SetOutput(logfile)

	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	dir := filepath.Dir(exe)
	realExe := filepath.Join(dir, common.AppName+".exe")

	cmd := exec.Command(realExe, os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}
}
