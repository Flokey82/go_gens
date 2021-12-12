package gamecs

type AgentMgr struct {
	entitiesByID map[int]*Agent
	entities     []*Agent
	nextID       int
}

func newAgentMgr() *AgentMgr {
	return &AgentMgr{
		entitiesByID: make(map[int]*Agent),
	}
}

func (m *AgentMgr) Entities() []*Agent {
	return m.entities
}

func (m *AgentMgr) NextID() int {
	id := m.nextID
	m.nextID++
	return id
}

func (m *AgentMgr) RegisterEntity(e *Agent) {
	m.entitiesByID[e.ID()] = e
	m.entities = append(m.entities, e)
}

func (m *AgentMgr) RemoveEntity(e *Agent) {
	delete(m.entitiesByID, e.ID())
	for i, en := range m.entities {
		if en == e {
			m.entities = append(m.entities[:i], m.entities[i+1:]...)
			return
		}
	}
}

func (m *AgentMgr) GetEntityFromID(id int) *Agent {
	return m.entitiesByID[id]
}

func (m *AgentMgr) Reset() {
	m.entitiesByID = make(map[int]*Agent)
	m.entities = nil
}
