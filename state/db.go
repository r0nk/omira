package state

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

func create_omira_db(filename string) *sql.DB {
	file, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	odb, _ := sql.Open("sqlite3", filename)
	create_task_statement = `CREATE TABLE tasks (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"name" TEXT,
	"due" DATETIME,
	"starting" DATETIME,
	"time_estimate" INT,
	"finished" DATETIME,
	"priority" REAL,
	"urgency" REAL,
	"recurrance" TEXT,
);`
	statement, err := odb.Prepare(create_task_statement) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	return odb
}

func load_task_db() {
	filename = "omira.db"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		odb = create_omira_db(filename)
	} else {
		odb, _ := sql.Open("sqlite3", filename)
	}
}
