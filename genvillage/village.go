// Package genvillage provides a very simple generator which attempts to generate a
// self-sustaining settlement economy given a number of buildings that consume and/or produce
// a number of resources.
package genvillage

import (
	"math/rand"
)

// Settlement represents a village with a number of buildings.
type Settlement struct {
	Buildings     []*Building // All buildings in the settlement.
	*Production               // Production totals for the settlement.
	*BuildingPool             // All known building types.
}

// NewSettlement returns a new settlement with the given BuildingPool.
func NewSettlement(p *BuildingPool) *Settlement {
	return &Settlement{
		BuildingPool: p,
		Production:   NewProduction(),
	}
}

// AddBuilding adds the given building to the settlement and updates the
// totals for each resource ID.
func (s *Settlement) AddBuilding(b *Building) {
	s.Buildings = append(s.Buildings, b)
	s.Update()
}

// Update rebuilds the totals for each resource ID required and/or provided
// by the local economy.
func (s *Settlement) Update() {
	s.Production.reset()
	for _, b := range s.Buildings {
		addToMap(s.Production.Requires, b.Requires)
		addToMap(s.Production.Provides, b.Provides)
	}
}

const maxAttempts = 1000

// Solve attempts to add known building types to the settlements until
// the local economy is self-sustaining.
func (s *Settlement) Solve() {
	// Here we'll keep track if the local economy is stable.
	var changed bool

	// Attempt to stabelize the local economy.
	for i := 0; i < maxAttempts; i++ {
		// Get the missing resources.
		req := s.GetMissing()
		if len(req) == 0 {
			break
		}

		// Try to find a buildings that provides the missing resources.
		for key, val := range req {
			bts := s.BuildingPool.Provides[key]
			if len(bts) == 0 {
				// We can't find any building that would provide the lacking resource.
				continue
			}

			// Add random buildings that provide the missing resource until
			// we have enough.
			for val > 0 {
				b := bts[rand.Intn(len(bts))]
				changed = true
				s.AddBuilding(b.NewBuilding())
				val -= b.Provides[key]
			}
		}

		// Nothing has changed in the last run, so return.
		if !changed {
			break
		}

		changed = false
	}
}
