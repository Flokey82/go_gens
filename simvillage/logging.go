package simvillage

import (
	"fmt"
	"log"
)

// HistoryManager manages game ticks and logging events
type HistoryManager struct {
	cal      *Calendar
	curr_log *Log
	logs     []*Log
	events   []string
}

func NewHistoryManager(game_timer *Calendar) *HistoryManager {
	return &HistoryManager{
		cal:      game_timer,
		curr_log: NewLog(game_timer.get_date()),
		logs:     nil,
	}
}

func (m *HistoryManager) Tick() []string {
	// Log old logs
	m.logs = append(m.logs, m.curr_log)

	// Log new log for logging new loggables
	m.curr_log = NewLog(m.cal.get_date())
	return nil
}

func (m *HistoryManager) add_event(event string) {
	m.events = append(m.events, event)
}

func (m *HistoryManager) add_events(events []string) {
	if len(events) == 0 {
		return
	}
	for _, event := range events {
		if event != "" {
			m.curr_log.add_event(0, event)
		}
	}
}
func (m *HistoryManager) list_todays_events() {
	m.curr_log.display_log()
}

// Log object represents a single tick in the village
type Log struct {
	date string

	world_events []string
	char_info    []string
	t1_events    []string
	t2_events    []string
	t3_events    []string
}

func NewLog(date string) *Log {
	return &Log{
		date: date,
	}
}

func (l *Log) add_event(tier int, text string) {
	if tier == 0 {
		l.world_events = append(l.world_events, text)
	}
	if tier == 1 {
		l.char_info = append(l.char_info, text)
	}
	if tier == 2 {
		l.t1_events = append(l.t1_events, "[1] "+text)
	}
	if tier == 3 {
		l.t2_events = append(l.t2_events, "[2] "+text)
	}
	if tier == 4 {
		l.t3_events = append(l.t3_events, "[3] "+text)
	}
}
func (l *Log) display_log() {
	var verbose_stack []string

	v := LOGGING_VERBOSITY
	if v > 0 {
		verbose_stack = append(verbose_stack, l.world_events...)
	}
	if v > 1 {
		verbose_stack = append(verbose_stack, l.char_info...)
	}
	if v > 2 {
		verbose_stack = append(verbose_stack, l.t1_events...)
	}
	if v > 3 {
		verbose_stack = append(verbose_stack, l.t2_events...)
	}
	if v > 4 {
		verbose_stack = append(verbose_stack, l.t3_events...)
	}
	log.Println(fmt.Sprintf("     ~~~ %s ~~~", l.date))

	for _, categories := range verbose_stack {
		log.Println(categories)
		//for _, items := range categories {
		//	log.Println(items)
		//}
	}
}
