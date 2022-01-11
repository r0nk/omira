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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		text.EnableColors()
		for _, string := range state.Finished_task_names {
			fmt.Printf("%s\n", text.Colors{text.FgHiBlack}.Sprintf("%s", string))
		}

		for _, string := range state.Unfinished_task_names {
			fmt.Printf("%s\n", string)
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
		fmt.Printf("\n")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
