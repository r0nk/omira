package state

import (
	"bufio"
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
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	today := time.Now()
	for scanner.Scan() {
		re := regexp.MustCompile("^" + today.Format("2006-01-02"))
		if re.MatchString(scanner.Text()) {
			field := strings.Fields(scanner.Text())
			Finished_task_names = append(Finished_task_names, field[1])
		}
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
