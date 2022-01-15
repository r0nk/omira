package state

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Date_to_path(date time.Time) string {
	return date.Format("calendar/2006/January/02")
}

func create_today() *bufio.Writer {
	path := Date_to_path(time.Now())
	os.MkdirAll(filepath.Dir(path), 0770)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	return bufio.NewWriter(file)
}

func read_todays_tasks() []string {
	f, err := os.Open(Date_to_path(time.Now()))
	var ret []string
	if err != nil {
		return ret
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	return ret
}

func check_deadline(t Task, current time.Time) {
	if t.Due.Before(current) {
		fmt.Printf("Task %s won't be completed in time\n", t.Name)
	}
}

func Schedule(working_hours float64) {
	writer := create_today()
	defer writer.Flush()
	var minutes_worked time.Duration
	var schedule []Task
	for _, t := range Tasks {
		if minutes_worked >= time.Duration(working_hours)*time.Hour {
			check_deadline(t, time.Now())
			continue
		}
		if t.Recurrance != "" && Should_recur(t.Recurrance, time.Now()) {
			copy_to_recurrance_directory(t)
		}
		if !t.Starting.Before(time.Now()) {
			continue
		}
		schedule = append(schedule, t)
		minutes_worked = minutes_worked + t.Time_estimate
		fmt.Fprintf(writer, "%s\n", t.Name)
	}
	if minutes_worked < (working_hours * 60) {
		fmt.Printf("Task queue underrun.")
	}
}
