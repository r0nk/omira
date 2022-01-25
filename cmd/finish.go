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
	Short: "Mark a task as complete",
	Long: `Finish a task, removing it from the task tree and adding it to the ledger.
For example:
	omira finish dance
`,
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
		err := os.Remove("tasks/" + task_to_remove.Name)
		if err != nil {
			fmt.Printf("Could not remove old task.\n")
			log.Fatal(err)
		}
		write_task_to_ledger(task_to_remove)
	},
}

func init() {
	rootCmd.AddCommand(finishCmd)
}
