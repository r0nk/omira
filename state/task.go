package state

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func Task_urgency(t Task) float64 {
	u := 0.0
	if time.Until(t.Starting) > 0 {
		return 0
	}
	u = 1000 - time.Until(t.Due).Hours() + float64((t.Priority * 10))
	if u < 0 {
		u = 0
	}
	return u
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

func Add_Task(t Task) {
	Tasks = append(Tasks, t)
	save_tasks()
}

//name time_estimate_in_minutes due_date finished_date
//project/thing 60
func parse_task(line string) Task {
	var t Task
	te := 0
	tokens := strings.Fields(line)
	//Partial definition, so we'll fill in the blanks
	t.Name = tokens[0]
	if len(tokens) == 1 {
		t.Time_estimate = time.Minute * time.Duration(60)
	} else {
		te, _ = strconv.Atoi(tokens[1])
		t.Time_estimate = time.Minute * time.Duration(te)
	}
	if len(tokens) == 2 {
		t.Due = time.Now()
	} else {
		t.Due = Get_date(tokens[2])
	}
	t.Urgency = Task_urgency(t)
	if len(tokens) == 4 {
		t.Finished = Get_date(tokens[3])
	}
	return t
}

func Load_Tasks() {
	file, err := os.Open("todo.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Tasks = append(Tasks, parse_task(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func save_tasks() {

	file, err := os.OpenFile("todo.txt", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	for _, t := range Tasks {
		if t.Finished.IsZero() {
			_, err = fmt.Fprintf(file, "%s	%.0f	%s\n", t.Name, t.Time_estimate.Minutes(), t.Due.Format("2006-01-02T15:04-07:00"))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			_, err = fmt.Fprintf(file, "%s	%.0f	%s	%s\n", t.Name, t.Time_estimate.Minutes(), t.Due.Format("2006-01-02T15:04-07:00"), t.Finished.Format("2006-01-02T15:04-07:00"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func Finish_Task(name string) {
	for i, t := range Tasks {
		if t.Name == name && t.Finished.IsZero() {
			Tasks[i].Finished = time.Now()
			break
		}
	}
	save_tasks()
}
