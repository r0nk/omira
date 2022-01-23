package state

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	ini "github.com/go-ini/ini"
	cp "github.com/otiai10/copy"
)

type Task struct {
	Name          string
	Due           time.Time `yaml:",flow"`
	Starting      time.Time
	Time_estimate time.Duration
	Priority      int
	Urgency       float64 `yaml:",omitempty"`
	Recurrance    string
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

func task_from_path(path string) Task {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	var t Task

	t.Name = strings.TrimPrefix(path, root_path)
	t.Due = Get_date(cfg.Section("").Key("due").String())
	s := cfg.Section("").Key("starting").String()
	if s != "" {
		t.Starting = Get_date(s)
	}
	t.Time_estimate = time.Minute * time.Duration(cfg.Section("").Key("time_est").MustInt())
	t.Priority, _ = cfg.Section("").Key("priority").Int()
	if t.Priority <= 0 {
		t.Priority = 1
	}
	t.Urgency = task_urgency(t)

	t.Recurrance = cfg.Section("").Key("recurrance").String()
	return t
}

func read_task(path string, d fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("Could not open task directory, did you add a task?\n")
		log.Fatal(err)
	}
	if d.Name()[0] == '.' {
		if d.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	}
	if d.IsDir() {
		return nil
	}
	Tasks = append(Tasks, task_from_path(path))
	return nil
}

func handle_recurring(path string, d fs.DirEntry, err error) error {
	if err != nil || d.IsDir() {
		return nil
	}
	t := task_from_path(path)

	if Should_recur(t.Recurrance, time.Now()) {
		fmt.Printf("")

		err := cp.Copy(root_path+t.Name, root_path+strings.TrimPrefix(t.Name, ".recurring"))
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func Insert_recurring_tasks() {
	filepath.WalkDir(root_path+".recurring", handle_recurring)
	Tasks = nil
	Load_Tasks()
}

func Load_Tasks() {
	root_path = "tasks/"
	filepath.WalkDir(root_path, read_task)
	sort.Slice(Tasks, func(p, q int) bool {
		return Tasks[p].Urgency > Tasks[q].Urgency
	})
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
