package cmd

import (
	"os"
	"os/signal"
	"syscall"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "omira",
	Short: "Omira is a personal task manager based on the unix philosophy.",
	Long: `
Omira is a personal task manager and scheduler that assigns tasks to be done every day.
It stores its information in an easily scriptable directory,
defined in its configuration file (default .omirarc.yaml).`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	state.Load_Tasks()
	signal.Ignore(syscall.SIGPIPE) // don't tell the user the pipe was closed early.
}
