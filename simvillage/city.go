package simvillage

import "fmt"

// CityManager manages the villages stockpiles and items.
type CityManager struct {
	name string

	pop        *PeopleManager
	population int

	wood  int
	stone int
	food  int

	log []string
}

func NewCityManager(people *PeopleManager) *CityManager {
	return &CityManager{
		name:       get_name(),
		pop:        people,
		population: len(people.people),
		wood:       20,
		stone:      20,
		food:       1000,
	}
}

func (c *CityManager) Tick() []string {
	c.population = len(c.pop.people)

	c.food_eaten()
	c.wood_burned()
	c.log_stats()

	cp_log := c.log
	c.log = nil
	return cp_log
}

func (c *CityManager) log_stats() {
	c.log = append(c.log, fmt.Sprintf("Pop:%d Wood:%d Stone:%d Food:%d\n", c.population, c.wood, c.stone, c.food))
}

func (c *CityManager) food_eaten() {

	sum_food := 0

	for _, p := range c.pop.people {
		if 0 < p.age && p.age < 4 {
			sum_food += 3
		}
		if 4 < p.age && p.age < 10 {
			sum_food += 4
		} else {
			sum_food += 5
		}
	}
	c.food -= sum_food

	if c.food < 0 {
		c.food = 0
		c.log = append(c.log, "Food stocks are empty! Bad times are ahead")
		for _, p := range c.pop.people {
			p.hunger -= 1
		}
	} else {
		c.log = append(c.log, fmt.Sprintf("The citizens eat %d food", sum_food))
		for _, p := range c.pop.people {
			if p.hunger < 11 {
				p.hunger += 1
			}
		}
	}
}
func (c *CityManager) wood_burned() {
	sum_burned_cook := 5
	sum_burned_warmth := 5

	c.wood -= (sum_burned_cook + sum_burned_warmth)

	if c.wood < 0 {
		c.wood = 0
		c.log = append(c.log, "Wood stocks are empty! Bad times are ahead")
	} else {
		c.log = append(c.log, fmt.Sprintf("The citizens burn %d wood to cook", sum_burned_cook))
		c.log = append(c.log, fmt.Sprintf("The citizens burn %d  wood to stay warm", sum_burned_warmth))
	}
}
