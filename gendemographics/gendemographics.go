// Package gendemographics is based on Medieval Demographics Made Easy by S. John Ross.
package gendemographics

import (
//"fmt"
//"log"
)

type Client struct {
	BT []*BusinessType
}

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
	s := NewSettlement(population)
	return s
}
