package cmd

import (
	"fmt"
	"path"
	"sort"
	"time"

	state "github.com/r0nk/omira/state"
	"github.com/spf13/cobra"
)

func projects_from_tasks() {
	sort.Slice(state.Tasks, func(i, j int) bool { return state.Tasks[i].Name < state.Tasks[j].Name })

	pmap := make(map[string][]state.Task)

	for _, t := range state.Tasks {
		name := path.Dir(t.Name)
		for name != "." {
			pmap[name] = append(pmap[name], t)
			name = path.Dir(name)
		}
	}
	for key, tasks := range pmap {
		var duration time.Duration
		var completed time.Duration
		var velocity time.Duration
		var left time.Duration
		recent := time.Now().AddDate(0, 0, -7) //1 week ago
		for _, task := range tasks {
			if !task.Finished.IsZero() {
				if task.Finished.After(recent) {
					velocity += task.Time_estimate
				}
				completed += task.Time_estimate
			}
			duration += task.Time_estimate
		}

		velocity /= 7
		if duration == completed { //dont print finished projects
			continue
		}
		left = duration - completed
		velocity_percentage := 100.0 * float32(velocity.Minutes()/duration.Minutes())
		//		fmt.Printf("velocity_percentage: %f\n", velocity_percentage)

		fmt.Printf("%s %s %0.1f%% +%0.1f%% ", key, left, 100*(completed.Hours()/duration.Hours()), velocity_percentage)
		if velocity == 0.0 {
			fmt.Printf("never")
		} else {
			fmt.Printf("%0.0f days", float32(left/velocity))
		}

		fmt.Printf("\n")
	}
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List projects and their completion",
	Long: `List projects, which are sets of tasks with the same prefixes.
Also lists the hours remaining and completion precentage.`,
	Run: func(cmd *cobra.Command, args []string) {
		projects_from_tasks()
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
