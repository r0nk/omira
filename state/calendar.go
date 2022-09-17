package state

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

func Discipline(t time.Time) float64 {
	d := t.Format("2006-01-02")
	x := len(Load_task_db("select * from tasks where strftime(\"%Y-%m-%d\",scheduled) == \"" + d + "\" and status != 'FINISHED'"))
	y := len(Load_task_db("select * from tasks where strftime(\"%Y-%m-%d\",scheduled) == \"" + d + "\" and status == 'FINISHED'"))

	if x+y == 0 {
		return 0
	}
	return 100.0 * (float64(y) / float64(x+y))
}

func set_scheduled_time(t Task) {
	if t.Due.Before(midnight_tonight()) {
		t.Scheduled = time.Now()
	}
	filename := "omira.db"

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Fatal("No omira.db file found cannot apply changes.")
	}

	odb, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer odb.Close()

	statement, err := odb.Prepare("UPDATE tasks SET scheduled = (datetime('now')) WHERE name = ? AND status != 'FINISHED'")
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec(t.Name)
}

func Schedule(working_hours float64) []Task {
	var minutes_worked time.Duration
	var total_urgency float64
	//	Insert_recurring_tasks()

	ks := knapsack.Get01Solution(uint64(working_hours*60), Tasks, func(t *Task) uint64 {
		ret := uint64(1)
		if !t.Starting.Before(time.Now()) {
			ret += 999999
		}
		ret += uint64(t.Time_estimate.Minutes())
		fmt.Printf("name: \"%s\" cost: %d\n", t.Name, ret)
		return ret
	}, func(t *Task) uint64 {
		ret := uint64(t.Urgency) + 1
		fmt.Printf("urgency: %d\n", ret)
		return ret
	})

	for _, t := range Tasks {
		if minutes_worked+t.Time_estimate >= time.Duration(working_hours)*time.Hour {
			check_deadline(t, time.Now())
		}
	}
	for _, t := range ks {
		minutes_worked += t.Time_estimate
		total_urgency += t.Urgency
		set_scheduled_time(t)
	}
	if minutes_worked < (time.Duration(working_hours) * 60) {
		fmt.Printf("Task queue underrun.\n")
	}
	if minutes_worked == 0 {
		fmt.Printf("Task queue empty.\n")
	}
	fmt.Printf("Scheduled %d tasks with total urgency %.0f to do in %0.1f/%0.1f hours\n", len(ks), total_urgency, minutes_worked.Hours(), working_hours)

	return ks
}
