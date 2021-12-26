package gamecs

import "log"

type Manager struct {
	entitiesByID  map[int]*Agent
	entities      []*Agent
	itemsByID     map[int]*Item
	items         []*Item
	locationsByID map[int]*Location
	locations     []*Location
	nextID        int
}

func newManager() *Manager {
	return &Manager{
		entitiesByID:  make(map[int]*Agent),
		itemsByID:     make(map[int]*Item),
		locationsByID: make(map[int]*Location),
	}
}

func (m *Manager) Locations() []*Location {
	return m.locations
}

func (m *Manager) RegisterLocation(loc *Location) {
	m.locationsByID[loc.ID()] = loc
	m.locations = append(m.locations, loc)
}

func (m *Manager) Items() []*Item {
	return m.items
}

func (m *Manager) RegisterItem(it *Item) {
	m.itemsByID[it.ID()] = it
	m.items = append(m.items, it)
}

func (m *Manager) RemoveItem(it *Item) {
	delete(m.itemsByID, it.ID())
	for i, in := range m.items {
		if in == it {
			m.items = append(m.items[:i], m.items[i+1:]...)
			if it.Location != LocWorld {
				m.GetEntityFromID(it.LocationID).CInventory.RemoveID(it.id)
			} else {
				log.Println("removed world item!!!!")
			}
			return
		}
	}
}

func (m *Manager) Entities() []*Agent {
	return m.entities
}

func (m *Manager) NextID() int {
	id := m.nextID
	m.nextID++
	return id
}

func (m *Manager) RegisterEntity(e *Agent) {
	m.entitiesByID[e.ID()] = e
	m.entities = append(m.entities, e)
}

func (m *Manager) RemoveEntity(e *Agent) {
	delete(m.entitiesByID, e.ID())
	for i, en := range m.entities {
		if en == e {
			m.entities = append(m.entities[:i], m.entities[i+1:]...)
			return
		}
	}
}

func (m *Manager) GetEntityFromID(id int) *Agent {
	return m.entitiesByID[id]
}

func (m *Manager) Reset() {
	m.entitiesByID = make(map[int]*Agent)
	m.entities = nil
}
