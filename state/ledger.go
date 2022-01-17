package state

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var Discipline float64
var Finished_task_names []string
var Unfinished_task_names []string

func read_omira_ledger(path string) {
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
	todays_tasks := read_todays_tasks()
	for _, t := range todays_tasks {
		found := false
		for _, f := range Finished_task_names {
			if f == t {
				found = true
				break
			}
		}
		if !found {
			Unfinished_task_names = append(Unfinished_task_names, t)
		}
	}
	Discipline = 100 * (1.0 - (float64(len(Unfinished_task_names)) / float64(len(todays_tasks))))
}
