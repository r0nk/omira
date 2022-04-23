package cmd

import (
	"fmt"

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
		for _, t := range state.Load_task_db("select * from tasks where scheduled = date('now');") {
			fmt.Printf("%s\n", t.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(leftCmd)
}
