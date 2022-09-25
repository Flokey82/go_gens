package gamecs

import "log"

// Manager is a rudimentary entity manager.
// NOTE: This sucks. I really need to re-think and
// re-write this.
type Manager struct {
	entitiesByID  map[int]*Agent
	entities      []*Agent
	itemsByID     map[int]*Item
	items         []*Item
	locationsByID map[int]*Location
	locations     []*Location
	nextID        int
}

// newManager returns a new entity manager.
func newManager() *Manager {
	return &Manager{
		entitiesByID:  make(map[int]*Agent),
		itemsByID:     make(map[int]*Item),
		locationsByID: make(map[int]*Location),
	}
}

// NextID returns the next available unique identifier.
func (m *Manager) NextID() int {
	id := m.nextID
	m.nextID++
	return id
}

// Reset resets stuff.
func (m *Manager) Reset() {
	m.entitiesByID = make(map[int]*Agent)
	m.entities = nil
	m.itemsByID = make(map[int]*Item)
	m.items = nil
	m.locationsByID = make(map[int]*Location)
	m.locations = nil
}

// Locations returns all registered locations.
func (m *Manager) Locations() []*Location {
	return m.locations
}

// RegisterLocation registers a new location.
func (m *Manager) RegisterLocation(loc *Location) {
	m.locationsByID[loc.ID()] = loc
	m.locations = append(m.locations, loc)
}

// Items returns all registered items.
func (m *Manager) Items() []*Item {
	return m.items
}

// RegisterItem registers the given item in the manager.
func (m *Manager) RegisterItem(it *Item) {
	m.itemsByID[it.ID()] = it
	m.items = append(m.items, it)
}

// RemoveItem removes the given item from the world.
func (m *Manager) RemoveItem(it *Item) {
	delete(m.itemsByID, it.ID())
	for i, in := range m.items {
		if in == it {
			m.items = append(m.items[:i], m.items[i+1:]...)
			if it.Location != LocWorld {
				m.GetEntityFromID(it.LocationID).CompInventory.RemoveID(it.id)
			} else {
				log.Println("removed world item!!!!")
			}
			return
		}
	}
}

// Entities returns all registered entities.
func (m *Manager) Entities() []*Agent {
	return m.entities
}

// RegisterEntity registers the given agent in the manager.
// NOTE: This should be generic and not be typed as *Agent.
func (m *Manager) RegisterEntity(e *Agent) {
	m.entitiesByID[e.ID()] = e
	m.entities = append(m.entities, e)
}

// RemoveEntity removes the given agent from the manager.
// NOTE: This should be generic and not be typed as *Agent.
func (m *Manager) RemoveEntity(e *Agent) {
	delete(m.entitiesByID, e.ID())
	for i, en := range m.entities {
		if en == e {
			m.entities = append(m.entities[:i], m.entities[i+1:]...)
			return
		}
	}
}

// GetEntityFromID returns the agent with the given ID (if any).
// NOTE: This should be generic and not be typed as *Agent.
func (m *Manager) GetEntityFromID(id int) *Agent {
	return m.entitiesByID[id]
}
