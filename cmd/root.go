package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg_file string

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

func init_config() {
	if cfg_file != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg_file)
		fmt.Printf("%s\n", cfg_file)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		//also serach current path by default
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".omirarc")
	}

	err := viper.ReadInConfig()
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		viper.SafeWriteConfigAs(home + "/.omirarc.yaml")
	}

	/*
		state.Working_path = viper.GetString("working_path")
		if state.Working_path == "" {
			fmt.Printf("\nQuick setup for config file and omira directory: \n\n")
			fmt.Printf(" $ mkdir omira/ && echo working_path: $(pwd)/omira >> ~/.omirarc.yaml \n\n")
			log.Fatal("working_path not found, it must be added to the config file or passed as a -p argument\n")
		}
		os.Chdir(state.Working_path)
		os.MkdirAll("tasks/.recurring/", 0770)
		state.Load()
	*/
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfg_file, "config", "", "Config file (default is $HOME/.omirarc.yaml)")
	rootCmd.PersistentFlags().StringP("working_path", "p", "", "Directory to read tasks and setup calender days")
	viper.BindPFlag("working_path", rootCmd.PersistentFlags().Lookup("working_path"))
	//cobra.OnInitialize(init_config)

	state.Load_Tasks()
	signal.Ignore(syscall.SIGPIPE) // don't tell the user the pipe was closed early.
}
