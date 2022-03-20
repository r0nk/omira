/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/text"
	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

func unicode_bar_from_percentage(x float64) string {
	switch {
	case x > 100-((100/8)*1):
		return "\u2588"
	case x > 100-((100/8)*2):
		return "\u2587"
	case x > 100-((100/8)*3):
		return "\u2586"
	case x > 100-((100/8)*4):
		return "\u2585"
	case x > 100-((100/8)*5):
		return "\u2584"
	case x > 100-((100/8)*6):
		return "\u2583"
	case x > 100-((100/8)*7):
		return "\u2582"
	}
	return "\u2581"
}

func discipline_percentage_color(x float64) (int, error) {
	switch {
	case x == 100.0:
		return fmt.Printf("%s", text.Colors{text.FgHiCyan}.EscapeSeq())
	case x > 80.0:
		return fmt.Printf("%s", text.Colors{text.FgGreen}.EscapeSeq())
	case x > 60.0:
		return fmt.Printf("%s", text.Colors{text.FgYellow}.EscapeSeq())
	}
	return fmt.Printf("%s", text.Colors{text.FgRed}.EscapeSeq())
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Pretty print todays tasks and progress bar.",
	Long: `Pretty print todays tasks and a progress bar. The bar fills based on
task completion, and is colored based on how much time is left in the day.
Finished tasks are greyed out, and the unfinished tasks are organized by time estimates.`,
	Run: func(cmd *cobra.Command, args []string) {
		text.EnableColors()
		for _, t := range state.Finished_task_names {
			fmt.Printf("%s\n", text.Colors{text.FgHiBlack}.Sprintf("%s", t))
		}

		var last_minutes_value float64
		last_minutes_value = -1.0
		for _, t := range state.Unfinished_tasks {
			fmt.Printf("%s", text.Colors{text.FgYellow}.EscapeSeq())
			if t.Time_estimate.Minutes() != last_minutes_value {
				fmt.Printf("%2.0f ", t.Time_estimate.Minutes())
				last_minutes_value = t.Time_estimate.Minutes()
			} else {
				fmt.Printf("   ")
			}
			fmt.Printf("\x1b[0m")
			fmt.Printf("%s\n", t.Name)
		}
		fmt.Printf("%s", text.Colors{text.FgCyan}.EscapeSeq())
		day_percentage := float64(100.0 * ((float64(time.Now().Hour()) - 7) / 16))

		for i := 0; i < 100; i += 2 {
			if float64(i) > day_percentage {
				fmt.Printf("%s", text.Colors{text.FgYellow}.EscapeSeq())
			}
			if state.Discipline < float64(i) {
				fmt.Printf("%s", "░")
			} else {
				fmt.Printf("%s", "█")
			}
		}
		discipline_percentage_color(state.Discipline)
		fmt.Printf(" %0.1f\n", state.Discipline)
		fmt.Printf("\x1b[0m")
		fmt.Printf("%s\n", unicode_bar_from_percentage(state.Discipline))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
