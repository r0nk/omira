package cmd

import (
	"fmt"

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
		fmt.Printf("TODO print the tasks that fit for today")
	},
}

func init() {
	rootCmd.AddCommand(leftCmd)
}
