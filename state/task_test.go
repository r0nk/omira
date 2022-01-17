package state

import (
	"testing"
	"time"
)

func TestShould_recur(t *testing.T) {
	if !Should_recur("* * *", time.Now()) {
		t.Error("all recurring not working")
	}

	if !Should_recur("2 16 *", Get_date("2001-02-16T12:00-06:00")) {
		t.Error("specific recurring not working")
	}
}
