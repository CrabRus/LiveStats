package domain

import "time"

type PeriodStats struct {
	StreamID  string
	StartedAt time.Time
	Words     map[string]int
	UserStats map[string]int
}
