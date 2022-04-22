package state

func Load() {
	Load_Tasks()
	read_omira_ledger("omira.ledger")
	load_task_db("select * from tasks")
}
