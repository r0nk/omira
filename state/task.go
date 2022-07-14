package state

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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
var Recurring []Task

var root_path string

func task_urgency(t Task) float64 {
	if time.Until(t.Starting) > 0 {
		return 0
	}
	return 1 - time.Until(t.Due).Hours() + float64((t.Priority * 10))
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
	for _, r := range Recurring {
		if Should_recur(r.Recurrance, time.Now()) {
			Tasks = append(Tasks, r)
		}
	}
}

func Load_Tasks() {
	Tasks = Load_task_db("select * from tasks")
	Recurring = Load_task_db("select * from recurring")
}

func midnight_tonight() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func Add_Task(t Task) {
	if t.Due.Before(midnight_tonight()) {
		t.Scheduled = time.Now()
	}
	filename := "/home/r0nk/life/omira.db"
	var tasks []Task

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Fatal("No omira.db file found cannot apply changes.")
	}

	odb, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer odb.Close()

	statement, err := odb.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec("INSERT INTO tasks (name,due,starting,time_estimate,finished,scheduled,priority,urgency,recurrance,status,notes) VALUES (?,?,?,?,?,?,?,?,?,?,?,?);", t.name, t.due, t.starting, t.time_t.estimate, t.finished, t.scheduled, t.priority, t.urgency, t.recurrance, t.status, t.notes)
}

func Finish_Task(t Task) {
	if t.Due.Before(midnight_tonight()) {
		t.Scheduled = time.Now()
	}
	filename := "/home/r0nk/life/omira.db"
	var tasks []Task

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Fatal("No omira.db file found cannot apply changes.")
	}

	odb, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer odb.Close()

	statement, err := odb.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec("UPDATE tasks SET status = 'FINISHED' WHERE name = ? AND status != 'FINISHED' LIMIT 1;", t.name)
}

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
