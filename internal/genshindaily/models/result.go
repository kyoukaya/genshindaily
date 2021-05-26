package models

import (
	"fmt"
	"time"
)

type Result struct {
	UID           string
	Today         string
	DaysCheckedIn int64
	Start         time.Time
	VersionText   string
	Award         Award
	Status        CheckInStatus
}

func (r *Result) String() string {
	return fmt.Sprintf("# %s\nDaily Report for ID %s: %s\n"+
		"Days checked in: %d\n"+
		"Award: %s x %d\n"+
		"Took %.2fs", r.VersionText,
		r.UID, r.Status.String(),
		r.DaysCheckedIn,
		r.Award.Name, r.Award.Cnt,
		time.Since(r.Start).Seconds(),
	)
}
