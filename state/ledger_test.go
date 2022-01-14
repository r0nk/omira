package state

import (
	"testing"
)

func TestRead_omira_ledger(t *testing.T) {
	read_omira_ledger("../tests/test.ledger")
	if Finished_task_names[0] != "meditation" {
		t.Errorf("First task does not match in file %s != %s",
			Finished_task_names[0], "make_bed")
	}
}
