// Package aistate implements a minimalistic state machine.
package aistate

// StateMachine implements a simple state machine.
type StateMachine struct {
	Current            State
	Previous           State
	transitions        map[StateType][]*Transition
	currentTransitions []*Transition
	anyTransitions     []*Transition
}

// New returns a new state machine.
func New() *StateMachine {
	return &StateMachine{
		transitions: make(map[StateType][]*Transition),
	}
}

// SetState sets the current state of the state machine.
func (s *StateMachine) SetState(state State) {
	if s.Current == state {
		return
	}
	if s.Current != nil {
		s.Current.OnExit()
	}
	s.Previous = s.Current
	s.Current = state
	s.currentTransitions = s.transitions[state.Type()]
	s.Current.OnEnter()
}

// RevertToPreviousState sets the current state of the state machine to its previous state.
func (s *StateMachine) RevertToPreviousState() {
	s.SetState(s.Previous)
}

// Tick advances the state machine by 'delta' (time elapsed since last tick).
func (s *StateMachine) Tick(delta uint64) {
	if t := s.GetTransition(); t != nil {
		s.SetState(t.to())
	}
	if s.Current != nil {
		s.Current.Tick(delta)
	}
}

// AddTransition adds a new transition from a state to another when predicate returns true.
func (s *StateMachine) AddTransition(from, to State, predicate func() bool) {
	s.transitions[from.Type()] = append(s.transitions[from.Type()], &Transition{
		to:        func() State { return to },
		condition: predicate,
	})
}

// AddAnyTransition adds a new transition that does not depend on a prior state, but only on
// the return value of 'predicate'.
func (s *StateMachine) AddAnyTransition(to State, predicate func() bool) {
	s.anyTransitions = append(s.anyTransitions, &Transition{
		to:        func() State { return to },
		condition: predicate,
	})
}

// AddSelector adds a new transition from a state to another returned by to() when predicate returns true.
func (s *StateMachine) AddSelector(from State, to func() State, predicate func() bool) {
	s.transitions[from.Type()] = append(s.transitions[from.Type()], &Transition{
		to:        to,
		condition: predicate,
	})
}

// AddAnySelector adds a new transition that does not depend on a prior state, but only on
// the return value of 'predicate' and returns a state returned by to().
func (s *StateMachine) AddAnySelector(to func() State, predicate func() bool) {
	s.anyTransitions = append(s.anyTransitions, &Transition{
		to:        to,
		condition: predicate,
	})
}

// GetTransition returns the next valid transition.
func (s *StateMachine) GetTransition() *Transition {
	for _, t := range s.anyTransitions {
		if t.condition() {
			return t
		}
	}
	for _, t := range s.currentTransitions {
		if t.condition() {
			return t
		}
	}
	return nil
}

// StateType represents the ID of a state.
type StateType int

// State defines an interface for a current state.
type State interface {
	Type() StateType   // Returns the ID of the current state
	Tick(delta uint64) // Advances the state and provides the 'delta' time elapsed.
	OnEnter()          // Hooks executed when transitioning to this state
	OnExit()           // Hooks executed when transitioning from this state
}

// Transition represents a transition to a specific state.
// TODO: Consider two different transition types, one with to as a function and one
// with to as a specific state (if performance is impacted by to being a function)
type Transition struct {
	to        func() State
	condition func() bool
}
