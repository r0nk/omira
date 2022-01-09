/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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

	state.Working_path = viper.GetString("working_path")
	if state.Working_path == "" {
		log.Fatal("working_path not found, it must be added to the config file or passed as a -p argument\n")
	}
	os.Chdir(state.Working_path)
	state.Load()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfg_file, "config", "", "Config file (default is $HOME/.config/omirarc)")
	rootCmd.PersistentFlags().StringP("working_path", "p", "", "Directory to read tasks and setup calender days")
	viper.BindPFlag("working_path", rootCmd.PersistentFlags().Lookup("working_path"))
	cobra.OnInitialize(init_config)

	signal.Ignore(syscall.SIGPIPE) // don't tell the user the pipe was closed early.
}
