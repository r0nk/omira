package data

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	ini "github.com/go-ini/ini"
)

type Task struct {
	Name          string
	Due           time.Time
	starting      time.Time
	Time_estimate time.Duration
	priority      int
	Urgency       float64
}

var Tasks []Task

var root_path string

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

func task_urgency(t Task) float64 {
	if time.Until(t.starting) > 0 {
		return 0
	}
	return 1 - time.Until(t.Due).Hours() + float64((t.priority * 10))
}

func read_task(path string, d fs.DirEntry, err error) error {
	if err != nil {
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
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	var t Task

	t.Name = strings.TrimPrefix(path, root_path)
	t.Due = get_date(cfg.Section("").Key("due").String())
	s := cfg.Section("").Key("starting").String()
	if s == "" {
		t.starting = time.Now()
	} else {
		t.starting = get_date(s)
	}
	t.Time_estimate = time.Minute * time.Duration(cfg.Section("").Key("time_est").MustInt())
	t.priority, _ = cfg.Section("").Key("priority").Int()
	if t.priority <= 0 {
		t.priority = 1
	}
	t.Urgency = task_urgency(t)

	Tasks = append(Tasks, t)
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

func Load_Tasks(path string) {
	//	read_omira_ledger(path)
	root_path = path + "/tasks/"
	filepath.WalkDir(root_path, read_task)
	sort.Slice(Tasks, func(p, q int) bool {
		return Tasks[p].Urgency > Tasks[q].Urgency
	})
}
