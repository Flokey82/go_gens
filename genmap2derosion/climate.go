package genmap2derosion

import (
	opensimplex "github.com/ojrac/opensimplex-go"
	"log"
	"math/rand"
	"time"
)

func (w *World) genClimate() *World2 {
	// Generate climate. This is really suboptimal.
	now := time.Now()
	w2 := newWorld2(w.dim.X, w.dim.Y, w.seed)
	w2.generate(w.heightmap)
	w.ExportPng("b_image_terrain.png", w2.terrain.heightmap)

	for i := 0; i < 365; i++ {
		log.Println(i)
		w2.day++
		w2.climate.calcWind(w2.day)
		w2.climate.calcTempMap()
		w2.climate.calcHumidityMap()
		w2.climate.calcRainMap()
		rm := make([]float64, len(w2.climate.RainMap))

		for i, v := range w2.climate.CloudMap {
			rm[i] = w.heightmap[i]
			if v {
				rm[i] = 0.2
			}
		}
		for i, v := range w2.climate.RainMap {
			if v {
				rm[i] = 0.7
			}
		}
		rm[0] = 0
		rm[1] = 1
		w.storeGifFrame(rm, w2.terrain.heightmap, w2.terrain.heightmap)
	}
	log.Println(w2.climate.WindMap)
	log.Println(w2.climate.HumidityMap)
	log.Println(w2.climate.RainMap)
	log.Println(w2.climate.CloudMap)
	w.ExportPng("b_image_avterrain.png", w2.terrain.heightmap)
	w.ExportPng("b_image_avgrain.png", w2.climate.AvgRainMap)
	w.ExportPng("b_image_avgtemp.png", w2.climate.AvgTempMap)
	w.ExportPng("b_image_avgwind.png", w2.climate.AvgWindMap)
	w.ExportPng("b_image_avgcloud.png", w2.climate.AvgCloudMap)
	log.Println(time.Since(now))
	return w2
}

type Climate struct {
	perlin     opensimplex.Noise
	seed       int
	dimX, dimY int
	//Curent Climate Maps
	TempMap       []float64
	HumidityMap   []float64
	CloudMap      []bool
	RainMap       []bool
	WindMap       []float64
	WindDirection [2]float64 //from 0-1

	//Average Climate Maps
	AvgRainMap     []float64
	AvgWindMap     []float64
	AvgCloudMap    []float64
	AvgTempMap     []float64
	AvgHumidityMap []float64
	terrain        *Terrain
}

func NewClimate(dimX, dimY, day, seed int, terrain *Terrain) *Climate {
	idxSize := dimX * dimY
	c := &Climate{
		dimX:           dimX,
		dimY:           dimY,
		seed:           seed,
		TempMap:        make([]float64, idxSize),
		HumidityMap:    make([]float64, idxSize),
		CloudMap:       make([]bool, idxSize),
		RainMap:        make([]bool, idxSize),
		WindMap:        make([]float64, idxSize),
		AvgRainMap:     make([]float64, idxSize),
		AvgWindMap:     make([]float64, idxSize),
		AvgCloudMap:    make([]float64, idxSize),
		AvgTempMap:     make([]float64, idxSize),
		AvgHumidityMap: make([]float64, idxSize),
		terrain:        terrain,
	}
	c.init(day)
	return c
}

func (c *Climate) init(day int) {
	if c.perlin == nil {
		c.perlin = opensimplex.New(int64(c.seed))
	}
	c.WindDirection[0] = 1
	c.WindDirection[1] = 1
	c.calcWind(day)
	c.initTempMap()
	c.initHumidityMap()
	c.initRainMap()
	c.initCloudMap()
}

func (c *Climate) calcAverage() {
	// Climate Simulation over n years
	years := 1
	startDay := 0

	// Initiate average climate maps
	for i := range c.terrain.heightmap {
		// Start at 0
		c.AvgRainMap[i] = 0
		c.AvgWindMap[i] = 0
		c.AvgCloudMap[i] = 0
		c.AvgTempMap[i] = 0
		c.AvgHumidityMap[i] = 0
	}

	//Initiate Simulation at a starting point
	simulation := NewClimate(c.dimX, c.dimY, startDay, c.seed, c.terrain)

	//Simulate every day for n years
	for i, days := 0, years*365; i < days; i++ {
		//Calculate new Climate
		simulation.calcWind(i)
		simulation.calcTempMap()
		simulation.calcHumidityMap()
		simulation.calcRainMap()

		// Calculate moving average.
		for idx := range c.terrain.heightmap {
			// Average wind.
			c.AvgWindMap[idx] = calcMovingAverage(c.AvgWindMap[idx], simulation.WindMap[idx], i)
			//c.AvgWindMap[j][k] = (c.AvgWindMap[j][k]*float64(i) + simulation.WindMap[j][k]) / float64(i+1)
			// Average rain.
			if simulation.RainMap[idx] {
				c.AvgRainMap[idx] = calcMovingAverage(c.AvgRainMap[idx], 1, i)
				//c.AvgRainMap[j][k] = (c.AvgRainMap[j][k]*float64(i) + 1) / float64(i+1)
			} else {
				c.AvgRainMap[idx] = calcMovingAverage(c.AvgRainMap[idx], 0, i)
				//c.AvgRainMap[j][k] = (c.AvgRainMap[j][k] * float64(i)) / float64(i+1)
			}

			// Average cloud cover.
			if simulation.CloudMap[idx] {
				c.AvgCloudMap[idx] = calcMovingAverage(c.AvgCloudMap[idx], 1, i)
				//c.AvgCloudMap[j][k] = (c.AvgCloudMap[j][k]*float64(i) + 1) / float64(i+1)
			} else {
				c.AvgCloudMap[idx] = calcMovingAverage(c.AvgCloudMap[idx], 0, i)
				//c.AvgCloudMap[j][k] = (c.AvgCloudMap[j][k] * float64(i)) / float64(i+1)
			}

			// Average temperature.
			c.AvgTempMap[idx] = calcMovingAverage(c.AvgTempMap[idx], simulation.TempMap[idx], i)
			//c.AvgTempMap[j][k] = (c.AvgTempMap[j][k]*float64(i) + simulation.TempMap[j][k]) / float64(i+1)

			// Average humidity.
			c.AvgHumidityMap[idx] = calcMovingAverage(c.AvgHumidityMap[idx], simulation.HumidityMap[idx], i)
			//c.AvgHumidityMap[j][k] = (c.AvgHumidityMap[j][k]*float64(i) + simulation.HumidityMap[j][k]) / float64(i+1)

		}
	}
}

func calcMovingAverage(v, newv float64, i int) float64 {
	return (v*float64(i) + newv) / float64(i+1)
}

func (c *Climate) calcWind(day int) {
	//Perlin Noise Module
	//var perlin Perlin
	//perlin.SetOctaveCount(2)
	//perlin.SetFrequency(4)

	timeInterval := float64(day) / 365

	// Winddirection shifts every Day
	// One Dimensional Perlin Noise
	c.WindDirection[0] = (c.perlin.Eval2(timeInterval, float64(c.seed)))
	c.WindDirection[1] = (c.perlin.Eval2(timeInterval, timeInterval+float64(c.seed)))

	dx := c.dimX
	dy := c.dimY
	wdx := c.WindDirection[0]
	wdy := c.WindDirection[1]
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			idx := i*dy + j
			// Previous Tiles
			k := i + int(10*wdx)
			if k < 0 || k >= dx {
				k = i
			}
			l := j + int(10*wdy)
			if l < 0 || l >= dy {
				l = j
			}

			c.WindMap[idx] = 5 * (1 - (c.terrain.heightmap[idx]-c.terrain.heightmap[k*dy+l])/1000)
		}
	}
}

func (c *Climate) initHumidityMap() {
	// Calculate the Humidity Grid
	for i := range c.HumidityMap {
		//Sea Level Temperature
		//c.HumidityMap[i] = 0 //In Degrees Celsius

		//Humidty Increases for
		if c.terrain.heightmap[i] < 200 {
			//In Degrees Celsius
			c.HumidityMap[i] = 0.4
		} else {
			c.HumidityMap[i] = 0.2
		}
	}
}

func (c *Climate) calcHumidityMap() {
	oldHumidMap := make([]float64, c.dimX*c.dimY)
	// Copy humidity to old humidity map.
	copy(oldHumidMap, c.HumidityMap)

	dx := c.dimX
	dy := c.dimY
	wdx := c.WindDirection[0]
	wdy := c.WindDirection[1]
	for i := 1; i < dx-1; i++ {
		for j := 1; j < dy-1; j++ {
			idx := i*dy + j
			// Get new map index from Wind Direction

			// Indices of Previous Tile
			// Assumption: Wind Blows Despite Obstacles
			k := i + int(2*c.WindMap[idx]*wdx)
			if k < 0 || k >= dx {
				k = i
			}
			l := j + int(2*c.WindMap[idx]*wdy)
			if l < 0 || l >= c.dimY {
				l = j
			}

			// Transfer to New Tile
			c.HumidityMap[idx] = oldHumidMap[k*dy+l]

			//Average
			newHumidity := (c.HumidityMap[idx-dy-1] + c.HumidityMap[idx+dy-1] + c.HumidityMap[idx+dy+1] + c.HumidityMap[idx-dy+1] + c.HumidityMap[idx] + c.HumidityMap[idx+1] + c.HumidityMap[idx-1] + c.HumidityMap[idx+dy] + c.HumidityMap[idx-dy]) / 9

			//newHumidity := (c.HumidityMap[idx-dy-1] + c.HumidityMap[idx+dy-1] + c.HumidityMap[idx+dy+1] + c.HumidityMap[idx-dy+1]) / 4

			// We are over a body of water, temperature accelerates
			var addHumidity float64
			if !c.CloudMap[idx] {
				if c.terrain.heightmap[idx] <= 200 {
					addHumidity = 0.05 * c.TempMap[idx]
				} else {
					addHumidity = 0.01
				}
			}

			// Raining
			var addRain float64
			if c.RainMap[idx] {
				addRain = -(newHumidity) * 0.8
			}

			newHumidity = newHumidity + (newHumidity)*addRain + (1-newHumidity)*(addHumidity)
			if newHumidity > 1 {
				newHumidity = 1
			} else if newHumidity < 0 {
				newHumidity = 0
			}
			c.HumidityMap[idx] = newHumidity
		}
	}
}

func (c *Climate) initTempMap() {
	for i := range c.TempMap {
		// Add for Height
		if h := c.terrain.heightmap[i]; h > 200 {
			// In Degrees Celsius
			c.TempMap[i] = 1 - h/2000
		} else {
			// Sea Temperature
			c.TempMap[i] = 0.7 // In Degrees Celsius
		}
	}
}

func (c *Climate) calcTempMap() {
	oldTempMap := make([]float64, c.dimX*c.dimY)

	// Copy temperature to old temperature map.
	copy(oldTempMap, c.TempMap)
	dx := c.dimX
	dy := c.dimY
	wdx := c.WindDirection[0]
	wdy := c.WindDirection[1]
	for i := 1; i < dx-1; i++ {
		for j := 1; j < dy-1; j++ {
			idx := i*dy + j
			// Get new map index from Wind Direction
			// Indices of Previous Tile
			k := i + int(2*c.WindMap[idx]*wdx)
			if k < 0 || k >= dx {
				k = i
			}
			l := j + int(2*c.WindMap[idx]*wdy)
			if l < 0 || l >= dy {
				l = j
			}

			// Transfer to New Tile
			c.TempMap[idx] = oldTempMap[k*dy+l]

			// Average
			newTemp := (c.TempMap[idx-dy-1] + c.TempMap[idx+dy-1] + c.TempMap[idx+dy+1] + c.TempMap[idx-dy+1] + c.TempMap[idx]) / 5

			// Various Contributions to the TempMap
			// Rising Air Cools
			addCool := 0.5 * (c.WindMap[idx] - 5)

			// Sunlight on Surface
			var addSun float64
			if !c.CloudMap[idx] {
				addSun = (1 - c.terrain.heightmap[idx]/2000) * 0.008
			}

			// Rain reduces temperature
			var addRain float64
			if c.RainMap[idx] && newTemp > 0 {
				addRain = -0.01
			}

			// Add Contributions
			newTemp = newTemp + 0.8*(1-newTemp)*addSun + 0.6*(newTemp)*(addRain+addCool)
			if newTemp > 1 {
				newTemp = 1
			} else if newTemp < 0 {
				newTemp = 0
			}
			c.TempMap[idx] = newTemp
		}
	}
}

func (c *Climate) initCloudMap() {
	for i := range c.CloudMap {
		c.CloudMap[i] = false
	}
}

func (c *Climate) initRainMap() {
	for i := range c.RainMap {
		c.RainMap[i] = false
	}
}

func (c *Climate) calcRainMap() {
	oldCloudMap := make([]bool, c.dimX*c.dimY)
	copy(oldCloudMap, c.CloudMap)

	oldRainMap := make([]bool, c.dimX*c.dimY)
	copy(oldRainMap, c.RainMap)

	for i := range oldRainMap {
		c.CloudMap[i] = false
		c.RainMap[i] = false
	}

	dx := c.dimX
	dy := c.dimY
	wdx := c.WindDirection[0]
	wdy := c.WindDirection[1]
	for i := 1; i < dx-1; i++ {
		for j := 1; j < dy-1; j++ {
			idx := i*dy + j
			// Old Coordinates
			k := i + int(2*c.WindMap[idx]*wdx)
			if k < 0 || k >= dx {
				k = i
			}
			l := j + int(2*c.WindMap[idx]*wdy)
			if l < 0 || l >= dy {
				l = j
			}

			// Rain Condition
			if c.HumidityMap[idx] >= 0.35+0.5*c.TempMap[idx] {
				c.RainMap[idx] = true
				// Transfer to New Tile
				c.CloudMap[idx] = oldCloudMap[k*dy+l]
			} else if c.HumidityMap[idx] >= 0.3+0.3*c.TempMap[idx] {
				c.CloudMap[idx] = true
				// Transfer to New Tile
				c.RainMap[idx] = oldRainMap[k*dy+l]
			} else {
				c.CloudMap[idx] = false
				c.RainMap[idx] = false
			}
		}
	}
}

type Terrain struct {
	seed        int
	worldDepth  int
	worldHeight int
	worldWidth  int

	// Terrain Parameters
	heightmap []float64
	biomeMap  []int
}

func newTerrain(dimX, dimY int64, seed int) *Terrain {
	idxSize := dimX * dimY
	return &Terrain{
		seed:        seed,
		worldDepth:  4000,
		worldHeight: int(dimY),
		worldWidth:  int(dimX),
		heightmap:   make([]float64, idxSize),
		biomeMap:    make([]int, idxSize),
	}
}

func (t *Terrain) genBiome(climate *Climate) {
	// Determine the Surface Biome:
	// 0: Water
	// 1: Sandy Beach
	// 2: Gravel Beach
	// 3: Stone Beach Cliffs
	// 4: Wet Plains (Grassland)
	// 5: Dry Plains (Shrubland)
	// 6: Rocky Hills
	// 7: Tempererate Forest
	// 8: Boreal Forest
	// 9: Mountain Tundra
	// 10: Mountain Peak
	// Compare the Parameters and decide what kind of ground we have.
	for i := range t.heightmap {
		switch d := t.heightmap[i]; {
		case d <= 200:
			t.biomeMap[i] = 0 // 0: Water
		case d <= 204:
			t.biomeMap[i] = 1 // 1: Sandy Beach
		case d <= 210:
			t.biomeMap[i] = 2 // 2: Gravel Beach
		case d <= 220:
			t.biomeMap[i] = 3 // 3: Stony Beach Cliffs
		case d <= 600:
			if climate.AvgRainMap[i] >= 0.02 {
				t.biomeMap[i] = 4 // 4: Wet Plains (Grassland)
			} else {
				t.biomeMap[i] = 5 // 5: Dry Plains (Shrubland)
			}
		case d <= 1300:
			x := i / t.worldHeight
			y := i % t.worldHeight
			if climate.AvgRainMap[i] < 0.001 && x+rand.Int()%4-2 > 5 && x+rand.Int()%4-2 < 95 && y+rand.Int()%4-2 > 5 && y+rand.Int()%4-2 < 95 {
				t.biomeMap[i] = 6 //6: Rocky Hills
			} else if d <= 1100 {
				t.biomeMap[i] = 7 //7: Temperate Forest
			} else {
				t.biomeMap[i] = 8 //8: Boreal Forest
			}
		case d <= 1500:
			t.biomeMap[i] = 9
		default:
			t.biomeMap[i] = 10 //Otherwise just Temperate Forest
		}
	}
}

func (t *Terrain) erode(years int) {
	//Climate Simulation
	average := NewClimate(t.worldWidth, t.worldHeight, 0, t.seed, t)

	//Simulate the Years
	for yr := 0; yr < years; yr++ {
		log.Println(yr)
		// Initiate the Climate
		average.init(0)

		// Simulate 1 Year for Average Weather Conditions
		average.calcAverage()

		// Add Erosion of the Climate after 1 Year
		var erosion float64
		for i := range t.heightmap {
			erosion = (average.AvgRainMap[i] + 0.5*average.AvgWindMap[i])
			curH := t.heightmap[i]
			t.heightmap[i] = curH - 5*(curH/2000)*(1-curH/2000)*erosion
		}
	}
}

func (t *Terrain) genHeight() {
	//Perlin Noise Module

	//Global Depth Map is Fine, unaffected by rivers.
	perlin := opensimplex.New(int64(t.seed))
	//var perlin Perlin

	//perlin.SetOctaveCount(12)
	//perlin.SetFrequency(2)
	//perlin.SetPersistence(0.6)

	//Generate the Perlin Noise World Map
	for i := range t.heightmap {
		//Generate the Height Map with Perlin Noise
		x := float64(i/t.worldHeight) / float64(t.worldWidth)
		y := float64(i%t.worldHeight) / float64(t.worldHeight)
		t.heightmap[i] = (perlin.Eval2(x, y))/5 + 0.25

		//Multiply with the Height Factor
		t.heightmap[i] *= float64(t.worldDepth)
	}
	log.Println(t.heightmap)
}

/*
func (t *Terrain) genLocal(seed int, player Player) {
	//Perlin Noise Module
	perlin := opensimplex.New(seed)
	//var perlin Perlin

	//perlin.SetOctaveCount(12)
	//perlin.SetFrequency(2)
	//perlin.SetPersistence(0.6)

	//Generate the Perlin Noise World Map
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			//Generate the Height Map with Perlin Noise
			x := float64(player.xTotal-25+i) / 100000
			y := float64(player.yTotal-25+j) / 100000
			t.localMap[i][j] = (perlin.Eval2(x, y, seed))/5 + 0.25
			//Multiply with the Height Factor
			t.localMap[i][j] = t.localMap[i][j] * t.worldDepth
		}
	}
}
*/
type Vegetation struct {
}

/*
func (v *Vegetation) getTree(territory World, player Player, i, j int) bool {
	//Code to Calculate wether or not we have a tree
	 Ideally this generates a vegetation map, spitting out
	0 for nothing,
	1 from short grass,
	2 for shrub,
	3 for some herb
	4 for some bush
	5 for some flower
	6 for some tree
	and also gives a number for a variant (3-5 variants of everything per biome)
	every variant could then also have a texture variant if wanted
	For now it only spits out wether or not we have a tree, which it then draws
	We can one piece of vegetation per map
	You could also do this for other objects on the map
	(tents, rocks, other locations) and not place vegetation if there is something present


	//Perlin Noise Module
	perlin := opensimplex.New(seed)
	// var perlin Perlin

	// perlin.SetOctaveCount(20)
	// perlin.SetFrequency(1000)
	// perlin.SetPersistence(0.8)

	//Generate the Height Map with Perlin Noise
	x := float64(player.xTotal-25+i) / 100000
	y := float64(player.yTotal-25+j) / 100000

	//This is not an efficient tree generation method
	//But a reasonable distribution for a grassland area
	//srand(x + y)
	tree := int((1/(perlin.Eval2(x, y, territory.seed+1)+1))*rand.Float64()%5) / 4

	return tree > 0
}*/

type World2 struct {
	seed       int
	day        int
	climate    *Climate
	terrain    *Terrain
	vegetation *Vegetation
}

func newWorld2(dimX, dimY, seed int64) *World2 {
	t := newTerrain(dimX, dimY, int(seed))
	return &World2{
		seed:       int(seed),
		day:        0,
		climate:    NewClimate(int(dimX), int(dimY), 0, int(seed), t),
		terrain:    t,
		vegetation: new(Vegetation),
	}
}

/*
func (w *World2) changePos(e SDL_Event) {
	switch e.key.keysym.sym {
	case SDLK_UP:
		yview -= 50
	case SDLK_DOWN:
		yview += 50
	case SDLK_LEFT:
		xview -= 50
	case SDLK_RIGHT:
		xview += 50
	}
}*/

func (w *World2) generate(h []float64) {
	//Geography
	//Generate and save a heightmap for all Blocks, all Regions
	w.terrain.genHeight()
	for i := range w.terrain.heightmap {
		w.terrain.heightmap[i] = h[i]*4000 - 300
	}

	//Erode the Landscape based on iterative average climate
	//w.terrain.erode(100)

	//Calculate the climate system of the eroded landscape
	w.climate.init(w.day)
	w.climate.calcAverage()

	//Generate the Surface Composition
	w.terrain.genBiome(w.climate)
}
