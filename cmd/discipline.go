/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

// disciplineCmd represents the discipline command
var disciplineCmd = &cobra.Command{
	Use:   "discipline",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%f\n", state.Discipline)
	},
}

func init() {
	rootCmd.AddCommand(disciplineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disciplineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// disciplineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
