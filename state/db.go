package state

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func create_omira_db(filename string) *sql.DB {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	odb, err := sql.Open("sqlite3", filename)

	if err != nil {
		log.Fatal(err)
	}
	create_task_statement := `
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
); `
	statement, err := odb.Prepare(create_task_statement) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	return odb
}

func Load_task_db(query string) []Task {
	filename := "/home/r0nk/life/omira.db"
	var tasks []Task

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Fatal("No omira.db file found, returning empty task list")
		return tasks
	}

	odb, err := sql.Open("sqlite3", filename)

	if err != nil {
		log.Fatal(err)
	}
	//if err != nil {
	//	odb = create_omira_db(filename)
	//}

	rows, err := odb.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	//TODO FIXME this is ugly
	for rows.Next() {
		var t Task
		var due sql.NullTime
		var starting sql.NullTime
		var finished sql.NullTime
		var scheduled sql.NullTime
		var time_estimate sql.NullInt32
		var priority sql.NullInt32
		var urgency sql.NullFloat64
		var recurrance sql.NullString
		var status sql.NullString
		var notes sql.NullString
		err = rows.Scan(&t.id, &t.Name, &due, &starting, &time_estimate, &finished, &scheduled, &priority, &urgency, &recurrance, &status, &notes)
		t.Due = due.Time
		t.Starting = starting.Time
		t.Finished = finished.Time
		t.Scheduled = scheduled.Time
		t.Time_estimate = time.Minute * time.Duration(time_estimate.Int32)
		t.Priority = int(priority.Int32)
		t.Urgency = urgency.Float64
		t.Recurrance = recurrance.String
		t.Status = status.String
		t.Notes = notes.String

		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, t)
	}

	defer rows.Close()

	/*
		for t := range tasks {
			fmt.Printf("name: %s\n", tasks[t].Name)
		}
	*/

	return tasks
}
