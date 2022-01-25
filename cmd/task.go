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

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "List all tasks in the task tree.",
	Long: `List all tasks in the task tree.
The first column displays urgency,
the second the task name,
the third the time estimate.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, t := range state.Tasks {
			fmt.Printf("%.0f %s %.0fh\n", t.Urgency, t.Name, time.Until(t.Due).Hours())
		}
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
