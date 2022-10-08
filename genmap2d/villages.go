package genmap2d

import (
	"math"
	"sort"
)

type VillageScore struct {
	X, Y  int
	Score float64
}

// PlaceVillage calculates the suitability score for each point on the map and
// will add a new village that gets the optimal score while being as far as possible
// from other villages.
// TODO: Factor out suitability function.
func (m *Map) PlaceVillage() {
	d := math.Sqrt(float64(m.Width*m.Width + m.Height*m.Height))
	var scores []*VillageScore
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if m.Cells[m.GetIndex(x, y)] != TileIDGrass {
				continue
			}
			ns := &VillageScore{X: x, Y: y}
			for _, cr := range m.tilesInRadius(x, y, 20) {
				switch cr {
				case TileIDTree:
					ns.Score += 0.01
				case TileIDWater:
					ns.Score += 0.01
				}
			}
			for _, other := range m.Villages {
				ns.Score -= float64(0.02 / (float64(dist(other.X, other.Y, ns.X, ns.Y))/d + 1e-9))
			}
			scores = append(scores, ns)
		}
	}
	if len(scores) == 0 {
		return
	}
	sort.Slice(scores, func(a, b int) bool {
		return scores[a].Score > scores[b].Score
	})
	winner := scores[0]
	m.Villages = append(m.Villages, winner)
	m.Cells[m.GetIndex(winner.X, winner.Y)] = TileIDVillage
}

// dist calculates the distance between two points.
func dist(x1, y1, x2, y2 int) int {
	return int(math.Sqrt(float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))))
}
