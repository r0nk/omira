package state

import (
	"fmt"
	"log"
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

func check_deadline(t Task, current time.Time) {
	if t.Due.Before(current) {
		fmt.Printf("Task %s won't be completed in time\n", t.Name)
	}
}

func Discipline(t time.Time) float64 {
	fmt.Printf("TODO Discipline\n")
	return 100
}

func Schedule(working_hours float64) {
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
	}
	for _, v := range ks {
		t := Tasks[v]
		if !t.Starting.Before(time.Now()) {
			continue
		}
		minutes_worked += t.Time_estimate
		schedule = append(schedule, t)
		//TODO ADD to today here
	}
	if minutes_worked < (time.Duration(working_hours) * 60) {
		fmt.Printf("Task queue underrun.")
	}
}
