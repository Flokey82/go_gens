package genworldvoronoi

import (
	"container/heap"
	"log"
)

// NOTE: This is an alternative empire implementation where we expand based
// on city states.

type Empire2 struct {
	ID      int      // Region where the empire originates (capital)
	Capital *City    // Capital city
	Culture *Culture // Culture of the empire
	Cities  []*City  // Cities within the territory
}

func (m *Civ) GetEmpires2() []*Empire2 {
	numEmpires := m.NumTerritories
	if numEmpires > m.NumCityStates {
		numEmpires = m.NumCityStates
	}
	sortCities := make([]*City, numEmpires)
	copy(sortCities, m.Cities)

	// TODO: Use cities with high expansionism.
	/*
		sort.Slice(sortCities, func(i, j int) bool {
			exA := sortCities[i].Score
			exB := sortCities[j].Score

			cA := m.getRCulture(sortCities[i].ID)
			if cA != nil {
				exA *= cA.Expansionism
			}
			cB := m.getRCulture(sortCities[j].ID)
			if cB != nil {
				exB *= cB.Expansionism
			}
			return exA > exB
		})
	*/

	var queue territoryQueue
	heap.Init(&queue)

	terr := initRegionSlice(len(m.Cities))
	cityIDToIndex := make(map[int]int)
	cityIDToCity := make(map[int]*City)
	for i, c := range m.Cities {
		cityIDToIndex[c.ID] = i
		cityIDToCity[c.ID] = c
	}
	for _, c := range sortCities {
		cc := m.GetCulture(c.ID)
		if cc == nil {
			log.Println("City has no culture", c.Name)
			continue
		}
		terr[cityIDToIndex[c.ID]] = c.ID
		for _, r := range m.getTerritoryNeighbors(c.ID, m.RegionToCityState) {
			log.Println("Adding", r, "to queue for", c.ID)
			newdist := m.getCityScoreForexp(cityIDToCity[r])
			heap.Push(&queue, &queueRegionEntry{
				score: newdist,
				city:  c.ID,
				vx:    r,
			})
		}

		log.Printf("City %s has expansionism %f and score %f", c.Name, cc.Expansionism, c.Score)
	}
	// Extend territories until the queue is empty.
	for queue.Len() > 0 {
		u := heap.Pop(&queue).(*queueRegionEntry)
		if terr[cityIDToIndex[u.vx]] >= 0 {
			continue
		}
		terr[cityIDToIndex[u.vx]] = u.city
		for _, v := range m.getTerritoryNeighbors(u.vx, m.RegionToCityState) {
			if terr[cityIDToIndex[v]] >= 0 {
				continue
			}
			newdist := m.getCityScoreForexp(cityIDToCity[v])
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueRegionEntry{
				score: newdist / u.score,
				city:  u.city,
				vx:    v,
			})
		}
	}

	log.Println(terr)

	// Now overwrite the territories with the new territories.
	// For this we will have to copy the city states and
	// set new territories.

	copy(m.RegionToTerritory, m.RegionToCityState)
	for i, t := range m.RegionToTerritory {
		if tn := terr[cityIDToIndex[t]]; tn >= 0 {
			m.RegionToTerritory[i] = tn
		}
	}

	return nil
}

func (m *Civ) getCityScoreForexp(c *City) float64 {
	cc := m.GetCulture(c.ID)
	if cc == nil {
		return c.Score
	}
	return c.Score * cc.Expansionism
}
