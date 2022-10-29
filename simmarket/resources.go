package simmarket

// Resource represents a single kind of resource.
type Resource any

// Resources is a map of resource type to amount or price.
type Resources map[Resource]float64

// MergeIn merges the given resouce set with the current one.
func (r Resources) MergeIn(o Resources) {
	for resource, units := range o {
		r[resource] = r[resource] + units
	}
}

// Clone makes a copy of the resource set.
func (r Resources) Clone() Resources {
	returnValue := make(Resources)
	for resource, units := range r {
		returnValue[resource] = units
	}
	return returnValue
}

// Eq returns true if the two sets are equal.
func (r Resources) Eq(o Resources) bool {
	for resource, units := range r {
		if otherUnits, ok := o[resource]; ok {
			if otherUnits != units {
				return false
			}
		} else {
			return false
		}
	}
	for resource, otherUnits := range o {
		if units, ok := r[resource]; ok {
			if units != otherUnits {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
