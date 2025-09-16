package cmd

import (
	"log/slog"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/AtifChy/aiub-notice/internal/logger"
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
			logger.L().Error("parsing num flag", slog.String("error", err.Error()))
			os.Exit(1)
		}

		numsMap := make(map[int]struct{})
		for _, n := range nums {
			numsMap[n] = struct{}{}
		}

		seen, err := notice.LoadSeenNotices()
		if err != nil {
			logger.L().Error("loading seen notices", slog.String("error", err.Error()))
			os.Exit(1)
		}
		if len(seen) == 0 {
			logger.L().Warn("no notices have been fetched yet")
			return
		}

		notices, err := notice.GetCachedNotices()
		if err != nil {
			logger.L().Error("fetching cached notices", slog.String("error", err.Error()))
			os.Exit(1)
		}

		sort.Slice(notices, func(i int, j int) bool {
			return notices[i].Date.After(notices[j].Date)
		})

		if len(notices) == 0 {
			logger.L().Warn("no new notices found")
			return
		}

		for idx, n := range notices {
			if _, ok := numsMap[idx+1]; !ok {
				continue
			}
			if _, ok := seen[n.Link]; ok {
				toast.Show(n)
				logger.L().Info("triggered toast for notice", slog.String("title", n.Title), slog.String("link", n.Link))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
	lastCmd.Flags().IntSliceP("num", "n", []int{1}, "Number(s) of last notices to display")
}
