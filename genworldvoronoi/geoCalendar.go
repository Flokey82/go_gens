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

// GetDayOfYear returns the current day of the year.
func (c *Calendar) GetDayOfYear() int {
	return c.t.YearDay()
}

// Tick advances the calendar by one day.
func (c *Calendar) Tick() {
	c.t = c.t.Add(time.Hour * 24)
}

// TickYear advances the calendar by one year.
func (c *Calendar) TickYear() {
	c.t = c.t.AddDate(1, 0, 0)
}

// SetYear sets the year of the calendar.
func (c *Calendar) SetYear(year int64) {
	c.t = time.Date(int(year), 1, 1, 0, 0, 0, 0, time.UTC)
}

// GetYear returns the current year.
func (c *Calendar) GetYear() int64 {
	return int64(c.t.Year())
}

// GetYearProgress returns the progress of the current year in 0.0-1.0.
func (c *Calendar) GetYearProgress() float64 {
	return float64(c.t.YearDay()) / float64(365)
}
