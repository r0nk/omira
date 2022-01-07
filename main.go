package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	ini "github.com/go-ini/ini"
)

type task struct {
	name          string
	due           time.Time
	starting      time.Time
	time_estimate time.Duration
	priority      int
	urgency       float64
}

var tasks []task

//" I hate time strings " - every programmer who ever lived
func get_date(str string) time.Time {
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

func task_urgency(t task) float64 {
	if time.Until(t.starting) > 0 {
		return 0
	}
	return 1 - time.Until(t.due).Hours() + float64((t.priority * 10))
}

func read_task(path string, d fs.DirEntry, err error) error {
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
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	var t task

	t.name = filepath.Base(path)
	t.due = get_date(cfg.Section("").Key("due").String())
	s := cfg.Section("").Key("starting").String()
	if s == "" {
		t.starting = time.Now()
	} else {
		t.starting = get_date(s)
	}
	t.time_estimate = time.Minute * time.Duration(cfg.Section("").Key("time_est").MustInt())
	t.priority, _ = cfg.Section("").Key("priority").Int()
	if t.priority <= 0 {
		t.priority = 1
	}
	t.urgency = task_urgency(t)

	tasks = append(tasks, t)
	return nil
}

func read_omira_ledger(path string) {
	f, err := os.Open(path + "/omira.ledger")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	const max_size = 50
	b := make([]byte, max_size)
	f.Read(b)
	fmt.Printf("b:%s\n", b)
}

func load_state(path string) {
	read_omira_ledger(path)
	filepath.WalkDir(path+"/tasks", read_task)
	sort.Slice(tasks, func(p, q int) bool {
		return tasks[p].urgency > tasks[q].urgency
	})
}

func main() {
	load_state("/home/r0nk/life")
	for _, t := range tasks {
		fmt.Printf("%s urgency: %f\n", t.name, t.urgency)
	}
}
