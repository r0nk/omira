package state

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
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
	create_task_statement := `CREATE TABLE tasks (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"name" TEXT,
	"due" DATETIME,
	"starting" DATETIME,
	"time_estimate" INT,
	"finished" DATETIME,
	"priority" REAL,
	"urgency" REAL,
	"recurrance" TEXT,
	"status" TEXT
);`
	statement, err := odb.Prepare(create_task_statement) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	return odb
}

func load_task_db() {
	filename := "omira.db"
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	odb, err := sqlx.Open("sqlite3", filename)

	if err != nil {
		log.Fatal(err)
	}
	//if err != nil {
	//	odb = create_omira_db(filename)
	//}

	rows, err := odb.Query("select * from tasks;")

	var tasks []Task

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for t := range tasks {
		fmt.Printf("name: %s\n", tasks[t].Name)
	}
}
