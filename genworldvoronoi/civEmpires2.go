package genworldvoronoi

import (
	"container/heap"
	"log"
	"sort"
)

// NOTE: This is an alternative empire implementation where we expand based
// on city states.

type Empire2 struct {
	ID      int      // Region where the empire originates (capital)
	Capital *City    // Capital city
	Culture *Culture // Culture of the empire
	Cities  []*City  // Cities within the territory
}

func (m *Map) GetEmpires2() []*Empire2 {
	numEmpires := 10
	if numEmpires > m.NumCityStates {
		numEmpires = m.NumCityStates
	}
	sortCities := make([]*City, len(m.cities_r))
	copy(sortCities, m.cities_r)
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

	var queue territoryQueue
	heap.Init(&queue)

	terr := initRegionSlice(len(m.cities_r))
	cityIDToIndex := make(map[int]int)
	cityIDToCity := make(map[int]*City)
	for i, c := range m.cities_r {
		cityIDToIndex[c.ID] = i
		cityIDToCity[c.ID] = c
	}
	for _, c := range sortCities {
		cc := m.getRCulture(c.ID)
		if cc == nil {
			continue
		}
		terr[cityIDToIndex[c.ID]] = c.ID
		for _, r := range m.getRTerritoryNeighbors(c.ID, m.r_city) {
			log.Println("Adding", r, "to queue for", c.ID)
			if terr[cityIDToIndex[r]] >= 0 {
				continue
			}
			newdist := m.getCityScoreForexp(cityIDToCity[r])
			terr[cityIDToIndex[r]] = c.ID
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
		for _, v := range m.getRTerritoryNeighbors(u.vx, m.r_city) {
			if terr[cityIDToIndex[v]] >= 0 {
				continue
			}
			newdist := m.getCityScoreForexp(cityIDToCity[v])
			if newdist < 0 {
				continue
			}
			heap.Push(&queue, &queueRegionEntry{
				score: u.score + newdist,
				city:  u.city,
				vx:    v,
			})
		}
	}

	log.Println(terr)

	return nil
}

func (m *Map) getCityScoreForexp(c *City) float64 {
	cc := m.getRCulture(c.ID)
	if cc == nil {
		return c.Score
	}
	return c.Score * cc.Expansionism
}
