package simmarket

// Carrier is the interface for a carrier of resources.
// NOTE: A carrier is responsible for fulfilling orders.
// I don't know why this is a separate interface from the trader.
type Carrier interface {
	Buy(*Order, *Order, float64)
	Deliver(*Order, *Order, float64)
}

// Trader is the interface for an individual merchant / trader.
type Trader interface {
	Asks() []*Order
	Bids() []*Order
}

// StandardTrader implements the Trader interface.
type StandardTrader struct {
	Carrier
	asks []*Order
	bids []*Order
}

// NewStandardTrader returns a new StandardTrader using a given carrier.
func NewStandardTrader(carrier Carrier) *StandardTrader {
	return &StandardTrader{carrier, nil, nil}
}

// Ask creates a new ask for a number of units of a given resource for the given price.
func (a *StandardTrader) Ask(units float64, resource Resource, price float64) {
	a.asks = append(a.asks, NewOrder(a, resource, units, price))
}

func (a *StandardTrader) Bid(units float64, resource Resource, price float64) {
	a.bids = append(a.bids, NewOrder(a, resource, units, price))
}

// Asks returns all outstanding asks.
func (a *StandardTrader) Asks() []*Order {
	return a.asks
}

// Bids returns all outstanding bids.
func (a *StandardTrader) Bids() []*Order {
	return a.bids
}
