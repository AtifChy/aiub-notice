package cmd

import (
	"log"

	"github.com/AtifChy/aiub-notice/internal/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all fetched notices",
	Long:  `Display all fetched notices in an interactive table that allows navigation and opening notices.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(list.NewModel())
		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running program: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
