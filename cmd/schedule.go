package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

var working_hours float64

func create_today() *bufio.Writer {
	path := state.Date_to_path(time.Now())
	os.MkdirAll(filepath.Dir(path), 0770)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	return bufio.NewWriter(file)
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a number of tasks to be completed today.",
	Long: `Schedule a number of tasks to be completed today.
Example;
	omira schedule -w 5 #schedule 5 hours worth of tasks today
	`,
	Run: func(cmd *cobra.Command, args []string) {
		writer := create_today()
		defer writer.Flush()
		var minutes_worked time.Duration
		var schedule []state.Task
		for _, t := range state.Tasks {
			if t.Starting.Before(time.Now()) {
				continue
			}
			schedule = append(schedule, t)
			minutes_worked = minutes_worked + t.Time_estimate
			fmt.Printf("%s\n", t.Name)
			fmt.Fprintf(writer, "%s\n", t.Name)
			if minutes_worked >= time.Duration(working_hours)*time.Hour {
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().Float64VarP(&working_hours, "working_hours", "w", 8, "Set the number of hours the tasks will fill.")
}
