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

var date_string string

// disciplineCmd represents the discipline command
var disciplineCmd = &cobra.Command{
	Use:   "discipline",
	Short: "Print the percentage of tasks complete today.",
	Long:  `Output the percentage of tasks complete to day, to be used in scripts or conkys.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%f\n", state.Discipline(time.Now()))
	},
}

func init() {
	rootCmd.AddCommand(disciplineCmd)
	disciplineCmd.Flags().StringVarP(&date_string, "date", "d", "", "Set the date to get discipline for.")
}
