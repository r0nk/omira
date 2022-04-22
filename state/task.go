package state

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	cp "github.com/otiai10/copy"
)

/*
CREATE TABLE tasks (
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        due DATETIME,
        starting DATETIME,
        time_estimate INT,
        finished DATETIME,
        scheduled DATETIME,
        priority REAL,
        urgency REAL,
        recurrance TEXT,
        status TEXT,
        notes TEXT
);
CREATE TABLE sqlite_sequence(name,seq);
*/

type Task struct {
	id            int
	Name          string
	Due           time.Time
	Starting      time.Time
	Time_estimate time.Duration
	Finished      time.Time
	Scheduled     time.Time
	Priority      int
	Urgency       float64
	Recurrance    string
	Status        string
	Notes         string
}

var Tasks []Task

var root_path string

func task_urgency(t Task) float64 {
	if time.Until(t.Starting) > 0 {
		return 0
	}
	return 1 - time.Until(t.Due).Hours() + float64((t.Priority * 10))
}

func copy_to_recurrance_directory(t Task) {
	err := cp.Copy(root_path+t.Name, root_path+".recurring/"+t.Name)
	if err != nil {
		log.Fatal(err)
	}
}

//https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/* this reason this takes the at parameter is in case we're scheduling for the future. */
func Should_recur(rstr string, at time.Time) bool {
	now_words := strings.Fields(fmt.Sprintf("%d %d %d \n", int(at.Month()), at.Day(), int(at.Weekday())))
	recur_words := strings.Fields(rstr)
	if len(recur_words) < 3 {
		log.Fatal("Improper Recurrance String: " + rstr)
	}

	for i, _ := range now_words {
		split_recur := strings.Split(recur_words[i], ",")
		if recur_words[i] != "*" && !stringInSlice(now_words[i], split_recur) {
			return false
		}
	}
	return true
}

func Task_from_name(name string) Task {
	for _, t := range Tasks {
		if t.Name == name {
			return t
		}
	}
	log.Fatal("Could not find task with name: " + name)
	var t Task
	return t
}

func Insert_recurring_tasks() {
	fmt.Printf("")
	/*TODO
	recurring := load_task_db("select * from recurring")
	for r := range recurring {
		if Should_recur(r.Recurrance, time.Now()) {
			Tasks = append(Tasks, t)
		}
	}
	*/
}

func Load_Tasks() {
	Tasks = load_task_db("select * from tasks")
}

func midnight_tonight() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func Add_Task(t Task) {
	var path string
	path = "tasks/" + t.Name
	os.MkdirAll(filepath.Dir(path), 0770)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "due:\t%s\n", t.Due.Format("2006-01-02T15:04-07:00"))
	fmt.Fprintf(writer, "time_est:\t%.0f\n", t.Time_estimate.Minutes())

	if t.Due.Before(midnight_tonight()) {
		append_today(t.Name)
	}
}
