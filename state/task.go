package state

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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
	tud := 0.0 //if the time isn't specified, assume its due now
	if !t.Due.IsZero() {
		tud = time.Until(t.Due).Hours()
	}
	u = 1000 - tud + float64((t.Priority * 10))
	if u < 0 {
		u = 0
	}
	return u
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
		t.Due = time.Time{} //empty time
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
		s := scanner.Text()
		if len(s) > 0 {
			Tasks = append(Tasks, parse_task(scanner.Text()))
		}
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
	sort.Slice(Tasks, func(i int, j int) bool {
		if Tasks[j].Finished.IsZero() {
			return true
		}
		return Tasks[j].Finished.Before(Tasks[i].Finished) //j<i
	})
	defer file.Close()
	for _, t := range Tasks {
		if t.Finished.IsZero() {
			if t.Due.IsZero() {
				_, err = fmt.Fprintf(file, "%s	%.0f\n", t.Name, t.Time_estimate.Minutes())
				if err != nil {
					fmt.Println(err)
				}
			} else {
				_, err = fmt.Fprintf(file, "%s	%.0f	%s\n", t.Name, t.Time_estimate.Minutes(), t.Due.Format("2006-01-02T15:04-07:00"))
				if err != nil {
					fmt.Println(err)
				}
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
