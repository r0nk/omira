package cmd

import (
	"fmt"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var ratio bool
var days int

// disciplineCmd represents the discipline command
var disciplineCmd = &cobra.Command{
	Use:   "discipline",
	Short: "Print the percentage of tasks complete today.",
	Long:  `Output the percentage of tasks complete to day, to be used in scripts or conkys.`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		for ; days >= 0; days-- {
			date := now.AddDate(0, 0, -days)
			//			fmt.Printf(date.Format(time.RFC850))

			discipline := state.Discipline(date)
			if ratio {
				minute := float64(date.Minute())
				hour := float64(date.Hour()) + (minute / 60)
				day_percentage := float64(100.0 * (hour - 9) / 8)
				fmt.Printf("%f\n", discipline/day_percentage)
			} else {
				fmt.Printf("%f\n", discipline)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(disciplineCmd)
	disciplineCmd.Flags().BoolVarP(&ratio, "ratio", "r", false, "Get the ratio of discipline versus time worked for the day.")
	disciplineCmd.Flags().IntVarP(&days, "days", "d", 1, "Get the discipline counts for the last n days.")
}
