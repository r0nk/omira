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
}

func parse_task(line string) Task {
	var t Task
	te := 0
	tokens := strings.Fields(line)
	if len(tokens) < 3 {
		fmt.Printf("INVALID LINE: %s\n", line)
		return t
	}
	t.Name = tokens[0]
	te, _ = strconv.Atoi(tokens[1])
	t.Time_estimate = time.Minute * time.Duration(te)
	t.Due = Get_date(tokens[2])
	t.Urgency = Task_urgency(t)
	if len(tokens) > 3 {
		t.Finished = Get_date(tokens[3])
	}
	return t
}

func Load_Tasks() {
	file, err := os.Open("example.txt")
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

func Finish_Task(name string) {
	fmt.Printf("TODO finish task %s\n", name)
}
