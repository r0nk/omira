package state

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/johnmuirjr/go-knapsack"
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

func check_deadline(t Task, current time.Time) {
	if t.Due.Before(current) {
		fmt.Printf("Task %s won't be completed in time\n", t.Name)
	}
}

func Tasks_finished_on(t time.Time) []Task {
	year, month, day := t.Date()
	//count the number of hours worked for that day, divide by 8
	var ret []Task
	for _, task := range Tasks {
		y, m, d := task.Finished.Date()
		if (y == year) && (m == month) && (d == day) {
			ret = append(ret, task)
		}
	}
	return ret
}

//Return the percentage of tasks done on a given day
func Discipline(t time.Time) float64 {
	year, month, day := t.Date()
	//count the number of hours worked for that day, divide by 8
	ret := 0.0
	for _, task := range Tasks {
		y, m, d := task.Finished.Date()
		if (y == year) && (m == month) && (d == day) {
			ret += task.Time_estimate.Minutes()
		}
	}
	ret = 100 * (ret / (60 * 8))
	return ret
}

func Schedule(working_hours float64) []Task {
	var minutes_worked time.Duration
	var total_urgency float64
	//	Insert_recurring_tasks()

	ks := knapsack.Get01Solution(uint64(working_hours*60), Tasks, func(t *Task) uint64 {
		ret := uint64(1)
		y, _, _ := t.Finished.Date()
		finished := (y != 1)
		if !t.Starting.Before(time.Now()) || finished {
			ret += 999999
		}

		ret += uint64(t.Time_estimate.Minutes())
		//		fmt.Printf("name: \"%s\" cost: %d\n", t.Name, ret)
		return ret
	}, func(t *Task) uint64 {
		ret := uint64(t.Urgency) + 1
		//		fmt.Printf("urgency: %d\n", ret)
		return ret
	})

	sort.Slice(ks, func(i int, j int) bool {
		return ks[i].Time_estimate < ks[j].Time_estimate
	})

	for _, t := range Tasks {
		if minutes_worked+t.Time_estimate >= time.Duration(working_hours)*time.Hour {
			check_deadline(t, time.Now())
		}
	}
	for _, t := range ks {
		minutes_worked += t.Time_estimate
		total_urgency += t.Urgency
	}
	if minutes_worked < (time.Duration(working_hours) * 60) {
		fmt.Printf("Task queue underrun.\n")
	}
	if minutes_worked == 0 {
		fmt.Printf("Task queue empty.\n")
	}
	//	fmt.Printf("Scheduled %d/%d tasks with total urgency %.0f to do in %0.1f/%0.1f hours\n", len(ks), len(Tasks), total_urgency, minutes_worked.Hours(), working_hours)

	return ks
}
