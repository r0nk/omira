/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var working_hours float64

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a number of tasks to be completed today.",
	Long: `Schedule a number of tasks to be completed today.
Example;
	omira schedule -w 5 #schedule 5 hours worth of tasks today
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var minutes_worked time.Duration
		var schedule []state.Task
		for _, t := range state.Tasks {
			schedule = append(schedule, t)
			minutes_worked = minutes_worked + t.Time_estimate
			fmt.Printf("%s\n", t.Name)
			if minutes_worked >= time.Duration(working_hours)*time.Hour {
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	// Here you will define your flags and configuration settings.

	scheduleCmd.Flags().Float64VarP(&working_hours, "working_hours", "w", 8, "Set the number of hours the tasks will fill.")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scheduleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scheduleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
