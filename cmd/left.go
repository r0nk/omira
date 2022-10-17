package cmd

import (
	"fmt"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var leftCmd = &cobra.Command{
	Use:   "left",
	Short: "Print the tasks that are left to do today.",
	Long: `Print the tasks that are left to do today.
Example:
	omira finish $(omira left | fzf)
`,
	Run: func(cmd *cobra.Command, args []string) {
		worked_minutes := 0.0
		for _, t := range state.Tasks_finished_on(time.Now()) {
			worked_minutes += t.Time_estimate.Minutes()
		}
		for _, t := range state.Schedule(6 - (worked_minutes / 60)) {
			fmt.Printf("%s\n", t.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(leftCmd)
}
