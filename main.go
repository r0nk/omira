package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	ini "github.com/go-ini/ini"
)

type task struct {
	name          string
	due           time.Time
	starting      time.Time
	time_estimate time.Duration
	priority      int
	urgency       float32
}

var tasks []task
var task_lock sync.Mutex

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

func read_task(wg *sync.WaitGroup, path string, name string) {
	defer wg.Done()
	task_lock.Lock()
	defer task_lock.Unlock()
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	var t task
	t.name = name

	t.due = get_date(cfg.Section("").Key("due").String())

	tasks = append(tasks, t)
}

func read_task_dir(wg *sync.WaitGroup, path string) {
	defer wg.Done()
	taskfiles, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range taskfiles {
		f_path := path + "/" + f.Name()
		if f.Name()[0] == '.' {
			continue
		}
		wg.Add(1)
		if f.IsDir() {
			go read_task_dir(wg, path+"/"+f.Name())
		} else {
			go read_task(wg, f_path, f.Name())
		}
	}
}

func read_omira_ledger(wg *sync.WaitGroup, path string) {
	defer wg.Done()
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
	var wg sync.WaitGroup
	wg.Add(2)
	go read_omira_ledger(&wg, path)
	go read_task_dir(&wg, path+"/tasks")
	wg.Wait()
}

func main() {
	load_state("/home/r0nk/life")
	for _, t := range tasks {
		fmt.Printf("%s due: %s\n", t.name, t.due.String())
	}
}
