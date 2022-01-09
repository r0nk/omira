/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func write_task_to_ledger(task_to_remove state.Task) {
	file, err := os.OpenFile("omira.ledger", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	bytes, err := yaml.Marshal(task_to_remove)
	if err != nil {
		log.Fatal(err)
	}
	n := time.Now()
	fmt.Fprintf(writer, "%s\t%s\n", n.Format("2006-01-02T15:04-07:00"), task_to_remove.Name)
	s := strings.Replace(string(bytes), "\n", "\n\t", -1)
	s = strings.Replace(s, "name: "+task_to_remove.Name+"\n", "", -1)
	fmt.Fprintln(writer, s)
}

// finishCmd represents the finish command
var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}
		var task_to_remove state.Task
		for _, t := range state.Tasks {
			if t.Name == args[0] {
				task_to_remove = t
				break
			}
		}
		if task_to_remove.Name == "" {
			fmt.Printf("Could not find task;%s\n", args[0])
			return
		}
		write_task_to_ledger(task_to_remove)
		os.Remove("tasks/" + task_to_remove.Name)
	},
}

func init() {
	rootCmd.AddCommand(finishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// finishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// finishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
