package state

import (
	"time"
)

func Date_to_path(date time.Time) string {
	return date.Format("calendar/2006/January/02")
}
