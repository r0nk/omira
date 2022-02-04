package state

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

//" I hate time strings " - every programmer who ever lived
func Get_date(str string) time.Time {
	t, err := time.Parse("2006-01-02T15:04-07:00", str)
	if err != nil {
		t, err = time.Parse("15:04", str)
		if err != nil {
			log.Fatal(err)
		}
		n := time.Now()
		t = t.AddDate(int(n.Year()), int(n.Month())-1, int(n.Day())-1)
	}
	return t
}

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

func append_today(s string) {
	path := Date_to_path(time.Now())
	file, err := os.OpenFile(path, os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(s)
}

func read_todays_task_names() []string {
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
	Insert_recurring_tasks()

	ks := Knapsack(len(Tasks), int(working_hours*60), func(i int) int {
		return int(Tasks[i].Urgency)
	}, func(i int) int {
		return int(Tasks[i].Time_estimate.Minutes())
	})

	var schedule []Task
	for _, t := range Tasks {
		if minutes_worked+t.Time_estimate >= time.Duration(working_hours)*time.Hour {
			check_deadline(t, time.Now())
			continue
		}
		if t.Recurrance != "" && Should_recur(t.Recurrance, time.Now()) {
			copy_to_recurrance_directory(t)
		}
	}
	for _, v := range ks {
		t := Tasks[v]
		if !t.Starting.Before(time.Now()) {
			continue
		}
		minutes_worked += t.Time_estimate
		schedule = append(schedule, t)
		fmt.Fprintf(writer, "%s\n", t.Name)
	}
	if minutes_worked < (time.Duration(working_hours) * 60) {
		fmt.Printf("Task queue underrun.")
	}
}
