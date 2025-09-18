package main

import (
	"os"

	"github.com/AtifChy/aiub-notice/cmd/aiub-notice/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
