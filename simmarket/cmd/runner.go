package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/Flokey82/go_gens/simmarket"
)

func main() {
	e := NewEconomy()
	for i := 0; i < 20; i++ {
		e.AddTrader(NewTrader(RandomRecipe()))
	}

	for i := 0; i < 100000; i++ {
		e.Run()
		fmt.Println(e)
	}
}

type Economy struct {
	Market  *simmarket.Market
	Traders []*Trader
}

func NewEconomy() *Economy {
	return &Economy{
		Market: simmarket.NewMarket(),
	}
}

func (e *Economy) AddTrader(t *Trader) {
	t.ID = len(e.Traders)
	e.Traders = append(e.Traders, t)
	e.Market.Add(t)
}
func (e *Economy) Run() {
	// Run the market.
	e.Market.Trade()

	// Run the traders.
	for _, t := range e.Traders {
		t.Produce()
		if t.failedToProduce > 10 {
			log.Printf("Trader %s failed to produce for 10 turns, changing recipe", t.Name())
			// Assign a new random recipe.
			t.Recipe = RandomRecipe()
			// Give some start capital
			t.Money = 1000.0
			// Reset
			t.failedToProduce = 0
		}
	}
}

func (e *Economy) String() string {
	var s string
	for _, t := range e.Traders {
		s += t.String()
	}
	return s
}

type carrier struct {
	need simmarket.Resources

	inventory simmarket.Resources
}

type Resource string

const (
	ResourceFood  Resource = "food"
	ResourceWood  Resource = "wood"
	ResourceIron  Resource = "iron"
	ResourceGrain Resource = "grain"
	ResourceTools Resource = "tools"
)

type Recipe struct {
	Name string
	In   simmarket.Resources
	Out  simmarket.Resources
}

type Trader struct {
	ID              int
	Money           float64
	MoneyPrev       float64
	Recipe          Recipe
	buyPrices       simmarket.Resources
	sellPrices      simmarket.Resources
	inventory       simmarket.Resources
	failedToProduce int
}

func NewTrader(recipe Recipe) *Trader {
	return &Trader{
		Money:     1000.0,
		MoneyPrev: 1000.0,
		Recipe:    recipe,
		buyPrices: simmarket.Resources{
			ResourceFood:  3.0,
			ResourceWood:  3.0,
			ResourceIron:  3.0,
			ResourceGrain: 3.0,
			ResourceTools: 6.0,
		},
		sellPrices: simmarket.Resources{
			ResourceFood:  1.0,
			ResourceWood:  1.0,
			ResourceIron:  1.0,
			ResourceGrain: 1.0,
			ResourceTools: 3.0,
		},
		inventory: simmarket.Resources{
			ResourceFood:  2,
			ResourceWood:  2,
			ResourceIron:  2,
			ResourceGrain: 2,
			ResourceTools: 2,
		},
	}
}

func (t *Trader) Name() string {
	return fmt.Sprintf("%s (%d)", t.Recipe.Name, t.ID)
}

func (t *Trader) String() string {
	var s string
	s += fmt.Sprintf("%s Money: %.2f", t.Name(), t.Money)
	return s
}

func (t *Trader) Produce() {
	// Update prices
	t.CalculateCost()

	// Check if we have enough resources to produce.
	var canProduce bool
	for resource, amount := range t.Recipe.In {
		if t.inventory[resource] < amount {
			canProduce = false
			break
		}
		canProduce = true
	}

	// Produce if we can.
	// TODO: Limit production by stock.
	if canProduce {
		for resource, amount := range t.Recipe.In {
			t.inventory[resource] -= amount
		}
		for resource, amount := range t.Recipe.Out {
			t.inventory[resource] += amount
			log.Printf("Trader %s produced %.2f %s", t.Name(), amount, resource)
		}
		t.failedToProduce = 0
	} else {
		t.failedToProduce++
		log.Printf("Trader %s failed to produce %d times", t.Name(), t.failedToProduce)
	}
}

func (t *Trader) CalculateCost() {
	// Check what we need to buy, sum up the cost
	sumCost := 0.0
	for resource, amount := range t.Recipe.In {
		sumCost += amount * t.buyPrices[resource]
	}

	// Calculate how much we can sell for
	countOutput := 0.0
	for _, amount := range t.Recipe.Out {
		countOutput += amount
	}

	// Calculate the price per unit
	pricePerUnit := sumCost / countOutput
	pricePerUnit = pricePerUnit + 0.1

	// Update the sell prices
	// TODO: This should be the new minimum price.
	for resource := range t.Recipe.Out {
		// Only correct upwards.
		t.sellPrices[resource] = pricePerUnit // TODO: add markup
	}
}

func (t *Trader) Buy(bid, ask *simmarket.Order, price float64) {
	log.Printf("Trader %s bought %.3f %s (%s) for %.2f", t.Name(), ask.Units, ask.Resource, bid.Resource, price)
	price = (bid.Price + ask.Price) / 2.0

	// Update the inventory
	t.inventory[bid.Resource] += t.inventory[bid.Resource] + bid.Units

	// Update the money
	t.Money = t.Money - bid.Units*price

	// Update the buy prices
	t.buyPrices[bid.Resource] = (t.buyPrices[bid.Resource] + price) / 2

	// Now make sure that the seller delivers.
	ask.Carrier.Deliver(bid, ask, price)
}

func (t *Trader) Deliver(bid, ask *simmarket.Order, price float64) {
	// Update the inventory
	t.inventory[ask.Resource] -= ask.Units

	// Update the money
	t.Money = t.Money + ask.Units*price

	// Update the sell prices
	t.sellPrices[ask.Resource] = (t.sellPrices[ask.Resource] + price) / 2
}

func (t *Trader) Asks() []*simmarket.Order {
	// Check what we can sell.
	canSell := make(simmarket.Resources)
	for resource, amount := range t.inventory {
		// Make sure we don't sell anything we need.
		if _, ok := t.Recipe.In[resource]; ok {
			continue
		}
		// Check how much we can sell.
		amountToSell := amount
		if amountToSell > 0 {
			canSell[resource] = amountToSell
		}
	}

	// Create the orders
	var orders []*simmarket.Order
	for resource, amount := range canSell {
		orders = append(orders, simmarket.NewOrder(t, resource, amount, t.sellPrices[resource]))
	}
	return orders
}

func (t *Trader) Bids() []*simmarket.Order {
	// Check what we need to buy.
	need := make(simmarket.Resources)
	var needMoney float64
	for resource, amount := range t.Recipe.In {
		// Check how much we have already in our inventory.
		amountInInventory := t.inventory[resource]
		// Check how much we need to buy.
		amountToBuy := amount - amountInInventory
		if amountToBuy > 0 {
			need[resource] = amountToBuy
			needMoney += amountToBuy * t.buyPrices[resource]
		}
	}

	// Do we have enough money?
	multiplier := 1.0
	if t.Money < needMoney {
		multiplier = t.Money / needMoney
	}

	// Create the orders
	var orders []*simmarket.Order
	for resource, amount := range need {
		orders = append(orders, simmarket.NewOrder(t, resource, amount, t.buyPrices[resource]*multiplier))
	}
	return orders
}

var (
	RecipeFarmer = Recipe{
		Name: "Farmer",
		In: simmarket.Resources{
			ResourceFood:  1,
			ResourceTools: 0.1,
		},
		Out: simmarket.Resources{
			ResourceGrain: 5,
		},
	}

	RecipeMiner = Recipe{
		Name: "Miner",
		In: simmarket.Resources{
			ResourceFood:  1,
			ResourceTools: 0.1,
		},
		Out: simmarket.Resources{
			ResourceIron: 1,
		},
	}

	RecipeBaker = Recipe{
		Name: "Baker",
		In: simmarket.Resources{
			ResourceGrain: 2,
		},
		Out: simmarket.Resources{
			ResourceFood: 3,
		},
	}

	RecipeBlacksmith = Recipe{
		Name: "Blacksmith",
		In: simmarket.Resources{
			ResourceFood: 1,
			ResourceIron: 1,
			ResourceWood: 1,
		},
		Out: simmarket.Resources{
			ResourceTools: 0.5,
		},
	}

	RecipeWoodcutter = Recipe{
		Name: "Woodcutter",
		In: simmarket.Resources{
			ResourceFood:  1,
			ResourceTools: 0.1,
		},
		Out: simmarket.Resources{
			ResourceWood: 4,
		},
	}
)

func RandomRecipe() Recipe {
	recipes := []Recipe{
		RecipeFarmer,
		RecipeMiner,
		RecipeBaker,
		RecipeBlacksmith,
		RecipeWoodcutter,
	}
	return recipes[rand.Intn(len(recipes))]
}
