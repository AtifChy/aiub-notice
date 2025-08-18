package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/AtifChy/aiub-notice/internal/common"
	"github.com/AtifChy/aiub-notice/internal/notice"
	"github.com/AtifChy/aiub-notice/internal/toast"
	"github.com/spf13/cobra"
)

// lastCmd represents the last command
var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Display the last fetched notice",
	Long:  `This command retrieves and displays the last fetched notice from the AIUB Notice Fetcher service.`,
	Run: func(cmd *cobra.Command, args []string) {
		seen, err := notice.LoadSeenNotices()
		if err != nil {
			log.Fatalf("Error loading seen notices: %v", err)
		}
		if len(seen) == 0 {
			fmt.Println("No notices have been fetched yet.")
			return
		}

		notices, err := notice.LoadNoticesCache()
		if err != nil {
			log.Fatalf("Error fetching notices: %v", err)
		}

		sort.Slice(notices, func(i int, j int) bool {
			return notices[i].Date.After(notices[j].Date)
		})

		for _, n := range notices {
			if _, ok := seen[n.Link]; ok {
				toast.Show(common.AUMID, n)
				fmt.Println("Triggered toast for:", n.Title)
				return
			}
		}

		fmt.Println("No new notices found.")
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
}
