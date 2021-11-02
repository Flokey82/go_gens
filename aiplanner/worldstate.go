package aiplanner

// WorldState provides an interface for querying and modifying the world state.
type WorldState interface {
	Get(key string) bool
	Fork() WorldFork
	Contains(otherState PlannerState) bool
}

// WorldFork provides an interface for querying and modifying a forked version of the actual world state.
type WorldFork interface {
	WorldState
	Set(key string, val bool)
	Update(otherState PlannerState)
}

// realWorld can be a fancy querying tool based on game mechanics and whatnot.
type realWorld map[string]bool

// Get queries the given key and returns a boolean value for it.
// TODO: This should allow for more complex queries depending on the actor that is planning.
// ... for example: If I want to know if "enemyVisible" to ensure line of sight, I need to know
// who is asking that question so I can determine if actor 'a' can see his enemy/target actor 'b' etc.
func (w *realWorld) Get(key string) bool {
	return (*w)[key]
}

// Set wouldn't be available in a real implementation.
// TODO: Remove from code.
func (w *realWorld) Set(key string, val bool) {
	(*w)[key] = val
}

// Fork returns a modifiable WorldState that can reflect a hypothetical future world state.
func (w *realWorld) Fork() WorldFork {
	return &modWorld{
		rw:   w,
		mods: make(map[string]bool),
	}
}

// Contains returns true if the real world state satisfies the 'otherState' state.
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

// Get returns the value set for the given key.
// NOTE: Returns a hypothetical value if set previously during a simulated action,
// otherwise it will try to get the actual value from the "real world".
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

// Set the given key to the given value without modifying the actual world state.
func (w *modWorld) Set(key string, val bool) {
	w.mods[key] = val
}

// Contains returns true if the hypothetical world state satisfies the 'otherState' state.
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

// Fork returns an independent, modifiable copy of this hypothetical future world state.
func (w *modWorld) Fork() WorldFork {
	nw := &modWorld{
		rw:   w.rw,
		mods: make(map[string]bool),
	}
	for key, val := range w.mods {
		nw.mods[key] = val
	}
	return nw
}

// PlannerState is a sad leftover type that needs to be looked at.
type PlannerState map[string]bool

func (s *PlannerState) Set(key string, value bool) {
	(*s)[key] = value
}
