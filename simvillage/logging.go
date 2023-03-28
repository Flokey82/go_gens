package simvillage

import (
	"log"
)

// HistoryManager manages game ticks and logging events
type HistoryManager struct {
	cal     *Calendar
	currLog *Log
	logs    []*Log
	events  []string
}

func NewHistoryManager(gameTimer *Calendar) *HistoryManager {
	return &HistoryManager{
		cal:     gameTimer,
		currLog: NewLog(gameTimer.getDate()),
	}
}

func (m *HistoryManager) Tick() []string {
	// Log old logs
	m.logs = append(m.logs, m.currLog)

	// Log new log for logging new loggables
	m.currLog = NewLog(m.cal.getDate())
	return nil
}

func (m *HistoryManager) addEvent(event string) {
	m.events = append(m.events, event)
}

func (m *HistoryManager) addEvents(events []string) {
	if len(events) == 0 {
		return
	}
	for _, event := range events {
		if event != "" {
			m.currLog.addEvent(0, event)
		}
	}
}

func (m *HistoryManager) listTodaysEvents() {
	m.currLog.displayLog()
}

// Log object represents a single tick in the village
type Log struct {
	date        string
	worldEvents []string
	charInfo    []string
	t1Events    []string
	t2Events    []string
	t3Events    []string
}

func NewLog(date string) *Log {
	return &Log{
		date: date,
	}
}

func (l *Log) addEvent(tier int, text string) {
	if tier == 0 {
		l.worldEvents = append(l.worldEvents, text)
	}
	if tier == 1 {
		l.charInfo = append(l.charInfo, text)
	}
	if tier == 2 {
		l.t1Events = append(l.t1Events, "[1] "+text)
	}
	if tier == 3 {
		l.t2Events = append(l.t2Events, "[2] "+text)
	}
	if tier == 4 {
		l.t3Events = append(l.t3Events, "[3] "+text)
	}
}

func (l *Log) displayLog() {
	var verboseStack []string
	v := LOGGING_VERBOSITY
	if v > 0 {
		verboseStack = append(verboseStack, l.worldEvents...)
	}
	if v > 1 {
		verboseStack = append(verboseStack, l.charInfo...)
	}
	if v > 2 {
		verboseStack = append(verboseStack, l.t1Events...)
	}
	if v > 3 {
		verboseStack = append(verboseStack, l.t2Events...)
	}
	if v > 4 {
		verboseStack = append(verboseStack, l.t3Events...)
	}
	log.Printf("     ~~~ %s ~~~\n", l.date)
	for _, categories := range verboseStack {
		log.Println(categories)
		// for _, items := range categories {
		//	log.Println(items)
		// }
	}
}
