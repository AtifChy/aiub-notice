package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/list"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "show"},
	Short:   "List all fetched notices",
	Long:    `Display all fetched notices in an interactive table that allows navigation and opening notices.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(list.NewModel())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("running interactive list: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
