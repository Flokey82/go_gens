package simvillage

import "fmt"

// Calendar manages time and seasons
type Calendar struct {
	day       int
	month     int
	year      int
	season    string
	sum_ticks int
	logs      []string
}

func NewCalendar() *Calendar {
	return &Calendar{
		day:       1,
		month:     3,
		year:      1,
		season:    "Spring",
		sum_ticks: 0,
	}
}

func (c *Calendar) Tick() []string {
	c.sum_ticks++
	c.increment_date()
	cp_logs := c.logs
	c.logs = nil
	return cp_logs
}

func (c *Calendar) increment_date() {
	c.day++
	if c.day%10 == 0 {
		c.day = 1
		c.month += 1
		c.set_season()
		c.logs = append(c.logs, fmt.Sprintf("It is now month %d", c.month))
		c.logs = append(c.logs, fmt.Sprintf("It is now %s", c.season))

		if c.month%12 == 0 {
			c.year++
			c.month = 1
			c.day = 1
			c.logs = append(c.logs, fmt.Sprintf("The year changed %d", c.year))
		}
	}
}
func (c *Calendar) set_season() {
	if c.month == 12 || c.month == 1 {
		c.season = "Winter"
	} else if c.month >= 2 && c.month <= 5 {
		c.season = "Spring"
	} else if c.month >= 6 && c.month <= 8 {
		c.season = "Summer"
	} else if c.month >= 9 && c.month <= 11 {
		c.season = "Fall"
	}
}

func (c *Calendar) get_date() string {
	return fmt.Sprintf("%d/%d/%d %s", c.month, c.day, c.year, c.season)
}
