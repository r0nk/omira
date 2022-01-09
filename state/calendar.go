package state

import (
	"bufio"
	"log"
	"os"
	"time"
)

func Date_to_path(date time.Time) string {
	return date.Format("calendar/2006/January/02")
}

func read_todays_tasks() []string {
	f, err := os.Open(Date_to_path(time.Now()))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var ret []string
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	return ret
}
