/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var urgencyCmd = &cobra.Command{
	Use:   "urgency",
	Short: "Print the sum of urgency of unfinished tasks.",
	Long:  `Print the sum of urgency of unfinished tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%f\n", state.Urgency())
	},
}

func init() {
	rootCmd.AddCommand(urgencyCmd)
}
