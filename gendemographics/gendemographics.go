// Package gendemographics is based on Medieval Demographics Made Easy by S. John Ross.
package gendemographics

// Client is a demographic generating client thingy. :P
type Client struct {
	BT []*BusinessType // business types
}

// New returns a new client for generating businesses based on demographics.
func New() *Client {
	return &Client{
		BT: BusinessTypes,
	}
}

// NewNation returns a new nation of given size in square miles (sorry) and population.
func (c *Client) NewNation(size, density int) *Nation {
	n := NewNation(size, density)
	for _, sz := range GenSettlementPopulations(n.Population()) {
		n.Settlements = append(n.Settlements, c.NewSettlement(sz))
	}
	return n
}

// NewSettlement returns a new settlement with the given population size.
func (c *Client) NewSettlement(population int) *Settlement {
	return NewSettlement(population)
}
