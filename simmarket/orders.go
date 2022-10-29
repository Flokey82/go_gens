package simmarket

import "fmt"

// Order represents a request to trade an amount of a resource for a given price.
type Order struct {
	Carrier  Carrier  // Creator of the offer
	Resource Resource // Resource to trade
	Units    float64  // Number of units to trade
	Price    float64  // Proposed price
}

func NewOrder(carrier Carrier, resource Resource, units, price float64) *Order {
	return &Order{
		Carrier:  carrier,
		Resource: resource,
		Units:    units,
		Price:    price,
	}
}

// String returns a string representation of the order.
func (o *Order) String() string {
	return fmt.Sprint(o.Carrier, ":", o.Resource, ":", o.Units, "*", o.Price)
}

// Orders is a sortable slice of orders.
type Orders []*Order

func (o Orders) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
func (o Orders) Len() int {
	return len(o)
}
func (o Orders) Less(i, j int) bool {
	return o[i].Price > o[j].Price
}
