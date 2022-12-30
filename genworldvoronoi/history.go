package genworldvoronoi

import "fmt"

type History struct {
	*Calendar
	Events []*Event
}

func NewHistory(c *Calendar) *History {
	return &History{
		Calendar: c,
	}
}

const (
	ObjectTypeCity = iota
	ObjectTypeCityState
	ObjectTypeEmpire
	ObjectTypeCulture
	ObjectTypeReligion
	ObjectTypeRegion
	ObjectTypeMountain
	ObjectTypeRiver
	ObjectTypeLake
	ObjectTypeSea
	ObjectTypeVolcano
)

type ObjectReference struct {
	ID   int  // ID of the object that the event is about.
	Type byte // Type of the object that the event is about.
}

type Event struct {
	Year int64           // TODO: Allow different types of time.
	Type string          // Maybe use an enum?
	Msg  string          // Message that describes the event.
	ID   ObjectReference // Reference to the object that the event is about.
}

func (e *Event) String() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Msg)
}

func (h *History) GetEvents(id int, t byte) []*Event {
	var events []*Event
	for _, e := range h.Events {
		if e.ID.ID == id && e.ID.Type == t {
			events = append(events, e)
		}
	}
	return events
}

func (h *History) AddEvent(t string, msg string, id ObjectReference) {
	h.Events = append(h.Events, &Event{
		Year: h.GetYear(),
		Type: t,
		Msg:  msg,
		ID:   id,
	})
}

/*
Ok, here are my thoughts on the history system:
Each event is connected to a specific object (e.g. a city, a person, a tribe,
a nation, a regionetc.). The event has a type (e.g. "city founded", "city
destroyed", "city conquered", "city renamed", "city renamed", "volcano erupted",
"earthquake", etc.). Of course, each event might have additional data (e.g.
"city destroyed" might have a reference to what destroyed the city, "city
conquered" might have a reference to who conquered the city, etc.).

When we tick the world, we have in each sub-system (e.g. the geo system, the
city system, the tribe system, the nation system, etc.) a list of events that
happened in the last tick. We then add these events to the history system.

Right now we use a simple slice to store the events, but this might become inefficient
if we have a lot of events. We might want to use a more efficient data structure
that allows us to quickly find events that happened in a specific time range.
There is for example a one dimensional range tree that might be useful, if we
want to find all events that happened in a specific time range.

Cross references might be useful to quickly identify events that are related to
a specific object or to each other.
The downside is, that we need to keep track of those cross references, which takes
additional memory... but once we switch to some form of database, we can probably
just store the cross references there and not in momory.
*/
