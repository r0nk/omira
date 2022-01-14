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

	Tasks = append(Tasks, t)
	return nil
}

func Load_Tasks() {
	//	read_omira_ledger(path)
	root_path = "tasks/"
	filepath.WalkDir(root_path, read_task)
	sort.Slice(Tasks, func(p, q int) bool {
		return Tasks[p].Urgency > Tasks[q].Urgency
	})
}

func Add_Task(t Task) {
	var path string
	path = "tasks/" + t.Name
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "due:\t%s\n", t.Due.Format("2006-01-02T15:04-07:00"))
	fmt.Fprintf(writer, "time_est:\t%.0f\n", t.Time_estimate.Minutes())

	/*
		bytes, err := yaml.Marshal(t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bytes))
	*/
}
