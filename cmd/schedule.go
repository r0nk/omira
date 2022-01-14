package cmd

import (
	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var working_hours float64

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a number of tasks to be completed today.",
	Long: `Schedule a number of tasks to be completed today.
Example;
	omira schedule -w 5 #schedule 5 hours worth of tasks today
	`,
	Run: func(cmd *cobra.Command, args []string) {
		state.Schedule(working_hours)
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().Float64VarP(&working_hours, "working_hours", "w", 8, "Set the number of hours the tasks will fill.")
}
