/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

// finishCmd represents the finish command
var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Mark a task as complete",
	Long: `Finish a task, marking it as complete.
For example:
	omira finish dance
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}
		var task_to_remove state.Task
		for _, t := range state.Tasks {
			if t.Name == args[0] {
				task_to_remove = t
				break
			}
		}
		if task_to_remove.Name == "" {
			fmt.Printf("Could not find task;%s\n", args[0])
			return
		}
		state.Finish_task(args[0])
	},
}

func init() {
	rootCmd.AddCommand(finishCmd)
}
