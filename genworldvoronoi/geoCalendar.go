package genworldvoronoi

import "time"

// TODO: Create struct that tracks date and time.
type Calendar struct {
	t time.Time
}

func NewCalendar() *Calendar {
	return &Calendar{
		t: time.Unix(0, 0),
	}
}

// GetYearProgress returns the progress of the current year in 0.0-1.0.
func (c *Calendar) GetYearProgress() float64 {
	return float64(c.t.YearDay()) / float64(365)
}
