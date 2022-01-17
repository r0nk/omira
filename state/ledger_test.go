package state

import (
	"testing"
)

func TestRead_omira_ledger(t *testing.T) {
	read_omira_ledger("../tests/test.ledger")
	if len(Finished_task_names) == 0 {
		t.Error("Empty Finished_task_names, could not test.")
		return
	}
	if Finished_task_names[0] != "meditation" {
		t.Errorf("First task does not match in file %s != %s",
			Finished_task_names[0], "make_bed")
	}
}
