package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/notice"
	"github.com/AtifChy/aiub-notice/internal/toast"
)

// lastCmd represents the last command
var lastCmd = &cobra.Command{
	Use:     "last",
	Aliases: []string{"recent"},
	Short:   "Display the last fetched notice",
	Long: `This command retrieves and displays the last fetched notice from the AIUB Notice Fetcher service.
Examples:
	# trigger toast for the last fetched notice
	aiub-notice last
	
	# trigger toast for multiple notices, e.g., last 1st, 3rd, and 5th notices
	aiub-notice last -n 1,3,5
`,
	Run: func(cmd *cobra.Command, args []string) {
		nums, err := cmd.Flags().GetIntSlice("num")
		if err != nil {
			fmt.Println("Error parsing num flag:", err)
		}

		numsMap := make(map[int]struct{})
		for _, n := range nums {
			numsMap[n] = struct{}{}
		}

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

		if len(notices) == 0 {
			fmt.Println("No new notices found.")
			return
		}

		for idx, n := range notices {
			if _, ok := numsMap[idx+1]; !ok {
				continue
			}
			if _, ok := seen[n.Link]; ok {
				toast.Show(n)
				fmt.Println("Triggered toast for:", n.Title)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
	lastCmd.Flags().IntSliceP("num", "n", []int{1}, "Number(s) of last notices to display")
}
