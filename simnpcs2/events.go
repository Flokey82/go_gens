package simnpcs2

// Events is an event manager for communicating events like
// damage taken, a being killed, etc.
type Events struct {
	Events []*Event // Recent events (will be cleared after each update).
}

// newEvents returns a new Events system.
func newEvents() *Events {
	return &Events{
		Events: make([]*Event, 0, 100),
	}
}

func (e *Events) Add(eventType EventType, source, target Entity, data interface{}) {
	event := &Event{
		EventType: eventType,
		Source:    source,
		Target:    target,
		Data:      data,
	}

	if target != nil {
		if n, ok := target.(Notifiable); ok {
			n.Notify(event)
		}
	}

	e.Events = append(e.Events, event)
}

// Update updates the Events system.
func (e *Events) Update(delta float64) {
	e.Events = e.Events[:0]
}

// Event represents an event that happened in the world.
type Event struct {
	EventType EventType
	Source    Entity
	Target    Entity
	Data      interface{}
}

// EventType represents the type of an event.
type EventType int

const (
	// EventAttack is an event that is triggered when an attack is made.
	EventAttack EventType = iota
)

// EventAttackData represents the data for an EventAttack event.
type EventAttackData struct {
	Damage float64
}

// EventListener is an event listener that can be added to an Entity.
type EventListener struct {
	Events []*Event
}

func newEventListener() *EventListener {
	return &EventListener{
		Events: make([]*Event, 0, 100),
	}
}

// FindType finds an event of a specific type.
// This is a temporary solution until we have a better event system.
func (l *EventListener) FindType(eventType EventType) *Event {
	for _, e := range l.Events {
		if e.EventType == eventType {
			return e
		}
	}
	return nil
}

// Notify notifies the listener of an event.
func (l *EventListener) Notify(event *Event) {
	l.Events = append(l.Events, event)
}

// Update updates the EventListener system.
func (l *EventListener) Update(delta float64) {
	l.Events = l.Events[:0]
}

// Notifiable is an interface that can be implemented by entities
// that want to be notified of events.
type Notifiable interface {
	Notify(event *Event)
}
