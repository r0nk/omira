package cmd

import (
	"fmt"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var due_string string
var task_to_add state.Task
var time_estimate float64
var waitcat bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task.",
	Long: `Add a task.
Example:
	omira add -d "1 week" -t 15 -n cook_the_bacon`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && task_to_add.Name == "" {
			cmd.Help()
			return
		}
		if task_to_add.Name == "" {
			task_to_add.Name = args[0]
		}

		if due_string == "1 week" {
			task_to_add.Due = time.Now().AddDate(0, 0, 7)
		} else {
			task_to_add.Due = state.Get_date(due_string)
		}
		task_to_add.Time_estimate = time.Duration(time_estimate) * time.Minute

		state.Add_Task(task_to_add)
		if waitcat {
			fmt.Printf("waitcat")
			// Read input into a file to be passed to start command
			// Print out the output passed to finish
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&due_string, "due", "d", "1 week", "Set the due date for the task.")
	addCmd.Flags().Float64VarP(&time_estimate, "time_estimate", "t", 60, "Set the estimated time of the task in minutes.")
	addCmd.Flags().StringVarP(&task_to_add.Name, "name", "n", "", "The name for the task. ")
	addCmd.Flags().BoolVarP(&waitcat, "waitcat", "w", false, "Activate waitcat mode, the program will wait until the task is complete.")
}
