package state

import (
	"fmt"
	"testing"
	"time"
)

func generate_fake_tasks() []Task {
	var tasks []Task
	var t Task

	t.Name = "Dishes"
	t.Due = time.Now().AddDate(0, 0, 7)
	t.Time_estimate = time.Duration(60) * time.Minute
	tasks = append(tasks, t)

	t.Name = "Vaccum"
	t.Due = time.Now()
	t.Time_estimate = time.Duration(120) * time.Minute
	tasks = append(tasks, t)

	t.Name = "Dust"
	t.Due = time.Now()
	t.Time_estimate = time.Duration(10) * time.Minute
	tasks = append(tasks, t)

	t.Name = "Clean"
	t.Due = time.Now()
	t.Time_estimate = time.Duration(90) * time.Minute
	tasks = append(tasks, t)

	t.Name = "Wash"
	t.Due = time.Now()
	t.Time_estimate = time.Duration(103) * time.Minute
	tasks = append(tasks, t)

	return tasks
}

func TestSchedule(t *testing.T) {
	Tasks = generate_fake_tasks()
	hours := 4.0
	scheduled_tasks := Schedule(hours)
	var total_time time.Duration
	for _, t := range scheduled_tasks {
		fmt.Printf("t.Name:%s\n", t.Name)
		total_time += t.Time_estimate
	}
	//Did it schedule anything at all?
	if total_time.Minutes() == 0 {
		t.Errorf("scheduled 0 minutes of work")
	}
	//Did it schedule over the line?
	if total_time.Minutes() > hours*60 {
		t.Errorf("Scheduled too much work")
	}
	//Did it schedule a task that wasn't started yet?
	//Did it schedule much less then it should?
}
