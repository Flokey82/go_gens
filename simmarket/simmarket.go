// Package simmarket implements a simple market simulation.
// See: https://download.tuxfamily.org/moneta/documents/LARC-2010-03.pdf
// Heavily based on: https://github.com/zond/gomarket
package simmarket

import (
	"math/rand"
	"sort"
)

// Market represents a single market where resources are traded.
type Market struct {
	traders map[Trader]bool
	prices  Resources
}

// NewMarket returns a new Market struct.
func NewMarket() *Market {
	return &Market{
		traders: make(map[Trader]bool),
		prices:  make(Resources),
	}
}

// Add adds the given trader to the market.
func (m *Market) Add(t Trader) {
	m.traders[t] = true
}

// Del removes the given trader from the market.
func (m *Market) Del(t Trader) {
	delete(m.traders, t)
}

// Price returns the price for a given resource if any is known.
func (m *Market) Price(r Resource) (price float64, ok bool) {
	price, ok = m.prices[r]
	return
}

// Value returns the total value of the given resources.
func (m *Market) Value(resources Resources) float64 {
	value := 0.0
	for resource, units := range resources {
		if price, ok := m.Price(resource); ok {
			value = value + price*units
		} else {
			value = value + units
		}
	}
	return value
}

// tradeResources attempts to resolve all given asks and bids and returns
// the total actual price the resources have been traded for.
func (m *Market) tradeResource(asks, bids []*Order) float64 {
	satisfied_bids := make(map[*Order]*Order)
	var lastAskPrice float64
	var lastBidPrice float64
	for len(asks) > 0 && len(bids) > 0 {
		ask := asks[len(asks)-1] // lowest ask
		bid := bids[0]           // highest bid
		lastAskPrice = ask.Price
		lastBidPrice = bid.Price

		// Check if we have a price match.
		if bid.Price < ask.Price {
			break // No match, so we are done.
		}
		if ask.Units > bid.Units {
			partial_ask := NewOrder(ask.Carrier, ask.Resource, bid.Units, ask.Price)
			satisfied_bids[bid] = partial_ask
			bids = bids[1:]
			ask.Units = ask.Units - bid.Units
		} else if ask.Units < bid.Units {
			partial_bid := NewOrder(bid.Carrier, bid.Resource, ask.Units, bid.Price)
			satisfied_bids[partial_bid] = ask
			asks = asks[:len(asks)-1]
			bid.Units = bid.Units - ask.Units
		} else {
			satisfied_bids[bid] = ask
			asks = asks[:len(asks)-1]
			bids = bids[1:]
		}
	}

	// Calculate the average price.
	var actualPrice float64
	if len(satisfied_bids) > 0 {
		if len(asks) == 0 && len(bids) == 0 {
			actualPrice = (lastAskPrice + lastBidPrice) / 2.0
		} else if len(asks) == 0 {
			actualPrice = lastBidPrice
		} else if len(bids) == 0 {
			actualPrice = lastAskPrice
		} else {
			actualPrice = (lastAskPrice + lastBidPrice) / 2.0
		}
	} else {
		actualPrice = (lastAskPrice + lastBidPrice) / 2.0
	}

	// Resolve all satisfied bids and complete the transactions.
	// NOTE: The clearing price should be agreed on by the two parties and
	// the traders should be notified on the trading volume, min, max, and
	// the clearing price.
	for bid, ask := range satisfied_bids {
		bid.Carrier.Buy(bid, ask, actualPrice)
	}
	return actualPrice
}

// Trade runs the trading simulation for one step.
func (m *Market) Trade() {
	// Sum up all asks and bids by resources.
	sums := m.createSums()
	allAsks := sums.asks
	allBids := sums.bids
	askSums := sums.askSums
	bidSums := sums.bidSums
	resources := sums.resources

	// Now commence all trades for each resource.
	for resource := range resources {
		asks := allAsks[resource]
		bids := allBids[resource]

		// Shuffle to avoid that traders with identical prices always
		// end up in the same order after sorting.
		rand.Shuffle(len(asks), func(i, j int) { asks[i], asks[j] = asks[j], asks[i] })
		rand.Shuffle(len(bids), func(i, j int) { bids[i], bids[j] = bids[j], bids[i] })

		// Sort asks and bids by price.
		sort.Sort(Orders(asks))
		sort.Sort(Orders(bids))

		if askSums[resource] == 0 {
			m.prices[resource] = bids[0].Price
		} else if bidSums[resource] == 0 {
			m.prices[resource] = asks[len(asks)-1].Price
		} else {
			m.prices[resource] = m.tradeResource(asks, bids)
		}
	}
}

// Sums is a helper struct to collect all asks and bids by resources.
type Sums struct {
	asks      map[Resource][]*Order // all asks per resource
	bids      map[Resource][]*Order // all bids per resource
	askSums   map[Resource]float64  // sum of all asks per resource
	bidSums   map[Resource]float64  // sum of all bids per resource
	resources map[Resource]bool     // is there an order for this resource?
}

// newSums returns a new Sums struct.
func newSums() *Sums {
	return &Sums{
		asks:      make(map[Resource][]*Order),
		bids:      make(map[Resource][]*Order),
		askSums:   make(map[Resource]float64),
		bidSums:   make(map[Resource]float64),
		resources: make(map[Resource]bool),
	}
}

// createSums creates the sums of all bids and asks for each resource.
func (m *Market) createSums() *Sums {
	sums := newSums()
	for trader := range m.traders {
		// Sum up all asks.
		for _, ask := range trader.Asks() {
			sums.asks[ask.Resource] = append(sums.asks[ask.Resource], ask)
			sums.askSums[ask.Resource] += ask.Units
			sums.resources[ask.Resource] = true
		}

		// Sum up all bids.
		for _, bid := range trader.Bids() {
			sums.bids[bid.Resource] = append(sums.bids[bid.Resource], bid)
			sums.bidSums[bid.Resource] += bid.Units
			sums.resources[bid.Resource] = true
		}
	}
	return sums
}
