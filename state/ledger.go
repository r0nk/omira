package state

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var Discipline float64
var Finished_task_names []string
var Unfinished_tasks []Task

func read_omira_ledger(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile("\n" + time.Now().Format("2006-01-02") + "[^\n].*")
	for _, str := range re.FindAllString(string(data), -1) {
		field := strings.Fields(str)
		Finished_task_names = append(Finished_task_names, field[1])
	}
	todays_tasks := read_todays_task_names()
	for _, t := range todays_tasks {
		found := false
		for _, f := range Finished_task_names {
			if f == t {
				found = true
				break
			}
		}
		if !found {
			Unfinished_tasks = append(Unfinished_tasks, Task_from_name(t))
		}
	}
	sort.Slice(Unfinished_tasks, func(p, q int) bool {
		if Unfinished_tasks[p].Time_estimate == Unfinished_tasks[q].Time_estimate {
			return strings.Compare(Unfinished_tasks[p].Name, Unfinished_tasks[q].Name) == -1
		} else {
			return Unfinished_tasks[p].Time_estimate < Unfinished_tasks[q].Time_estimate
		}
	})
	Discipline = 100 * (1.0 - (float64(len(Unfinished_tasks)) / float64(len(todays_tasks))))
	return nil
}
