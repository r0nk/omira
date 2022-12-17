package cmd

import (
	"fmt"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var date_string string
var ratio bool

// disciplineCmd represents the discipline command
var disciplineCmd = &cobra.Command{
	Use:   "discipline",
	Short: "Print the percentage of tasks complete today.",
	Long:  `Output the percentage of tasks complete to day, to be used in scripts or conkys.`,
	Run: func(cmd *cobra.Command, args []string) {
		discipline := state.Discipline(time.Now())
		if ratio {
			minute := float64(time.Now().Minute())
			hour := float64(time.Now().Hour()) + (minute / 60)
			day_percentage := float64(100.0 * (hour - 9) / 8)
			fmt.Printf("%f\n", discipline/day_percentage)
		} else {
			fmt.Printf("%f\n", discipline)
		}
	},
}

func init() {
	rootCmd.AddCommand(disciplineCmd)
	disciplineCmd.Flags().StringVarP(&date_string, "date", "d", "", "Set the date to get discipline for.")
	disciplineCmd.Flags().BoolVarP(&ratio, "ratio", "r", false, "Get the ratio of discipline versus time worked for the day.")
}
