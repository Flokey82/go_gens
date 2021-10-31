package aiplanner

// WorldState provides an interface for querying and modifying the world state.
type WorldState interface {
	Get(key string) bool
	Set(key string, val bool)
	Fork() WorldState
	Update(otherState PlannerState)
	Contains(otherState PlannerState) bool
}

// realWorld can be a fancy querying tool based on game mechanics and whatnot.
type realWorld map[string]bool

func (w *realWorld) Get(key string) bool {
	return (*w)[key]
}

// Set wouldn't be available in a real implementation.
// TODO: Remove from code.
func (w *realWorld) Set(key string, val bool) {
	(*w)[key] = val
}

func (w *realWorld) Fork() WorldState {
	return &modWorld{
		rw:   w,
		mods: make(map[string]bool),
	}
}

func (w *realWorld) Contains(otherState PlannerState) bool {
	for key, value := range otherState {
		ourValue, ok := (*w)[key]
		if !ok || value != ourValue {
			return false
		}
	}
	return true
}

// modWorld represents the world state with an overlay that represents hypothetical changes to the world state.
type modWorld struct {
	rw   *realWorld      // rw contains the actual queryable world state.
	mods map[string]bool // mods contains modifications to the world state
}

func (w *modWorld) Get(key string) bool {
	if v, ok := w.mods[key]; ok {
		return v
	}

	// TODO: Cache result and somehow distribute cached result to other forks
	// that do not have yet set a value for this or introduce a global caching layer
	// in the realWorld or something. This will allow us to skip expensive calculations
	// after we have queried them once.
	return w.rw.Get(key)
}

func (w *modWorld) Set(key string, val bool) {
	w.mods[key] = val
}

func (w *modWorld) Contains(otherState PlannerState) bool {
	for key, value := range otherState {
		if ourValue, ok := w.mods[key]; ok {
			if value != ourValue {
				return false
			}
		} else if w.rw.Get(key) != value {
			return false
		}
	}
	return true
}

func (w *modWorld) Update(otherState PlannerState) {
	for key, value := range otherState {
		w.mods[key] = value
	}
}

func (w *modWorld) Fork() WorldState {
	nw := &modWorld{
		rw:   w.rw,
		mods: make(map[string]bool),
	}
	for key, val := range w.mods {
		nw.mods[key] = val
	}
	return nw
}

type PlannerState map[string]bool

func (s *PlannerState) Contains(otherState PlannerState) bool {
	for key, value := range otherState {
		ourValue, ok := (*s)[key]
		if !ok || value != ourValue {
			return false
		}
	}
	return true
}

func (s *PlannerState) Diff(otherState PlannerState) (newState PlannerState, hasDiff bool) {
	newState = make(PlannerState)
	for key, value := range otherState {
		ourValue, ok := (*s)[key]
		if !ok || value != ourValue {
			hasDiff = false
			newState[key] = value
		}
	}
	if hasDiff {
		return newState, hasDiff
	}
	return nil, false
}

func (s *PlannerState) Update(otherState PlannerState) {
	for key, value := range otherState {
		(*s)[key] = value
	}
}

func (s *PlannerState) Set(key string, value bool) {
	(*s)[key] = value
}
