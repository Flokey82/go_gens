package simvillage

import "fmt"

// CityManager manages the villages stockpiles and items.
type CityManager struct {
	pop        *PeopleManager
	name       string
	population int
	wood       int
	stone      int
	food       int
	log        []string
}

func NewCityManager(people *PeopleManager) *CityManager {
	return &CityManager{
		pop:        people,
		name:       getName(),
		population: len(people.people),
		wood:       20,
		stone:      20,
		food:       1000,
	}
}

func (c *CityManager) Tick() []string {
	c.population = len(c.pop.people)
	c.foodEaten()
	c.woodBurned()
	c.logStats()
	cp_log := c.log
	c.log = nil
	return cp_log
}

func (c *CityManager) logStats() {
	c.log = append(c.log, fmt.Sprintf("Pop:%d Wood:%d Stone:%d Food:%d\n", c.population, c.wood, c.stone, c.food))
}

func (c *CityManager) foodEaten() {
	sumFood := 0
	for _, p := range c.pop.people {
		if 0 < p.age && p.age < 4 {
			sumFood += 1
		} else if 4 < p.age && p.age < 10 {
			sumFood += 3
		} else {
			sumFood += 5
		}
	}
	c.food -= sumFood
	if c.food < 0 {
		c.food = 0
		c.log = append(c.log, "Food stocks are empty! Bad times are ahead")
		for _, p := range c.pop.people {
			p.hunger -= 1
		}
	} else {
		c.log = append(c.log, fmt.Sprintf("The citizens eat %d food", sumFood))
		for _, p := range c.pop.people {
			if p.hunger < 11 {
				p.hunger += 1
			}
		}
	}
}

func (c *CityManager) woodBurned() {
	sumBurnedCook := 5
	sumBurnedWarmth := 5
	c.wood -= (sumBurnedCook + sumBurnedWarmth)
	if c.wood < 0 {
		c.wood = 0
		c.log = append(c.log, "Wood stocks are empty! Bad times are ahead")
	} else {
		c.log = append(c.log, fmt.Sprintf("The citizens burn %d wood to cook", sumBurnedCook))
		c.log = append(c.log, fmt.Sprintf("The citizens burn %d wood to stay warm", sumBurnedWarmth))
	}
}
