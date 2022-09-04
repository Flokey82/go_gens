// Package gamedice is a pretty simple dice roller for games.
package gamedice

import "math/rand"

// Dice represents a dice.
type Dice int

// Various dice types.
const (
	D2   Dice = 2
	D4   Dice = 4
	D6   Dice = 6
	D8   Dice = 8
	D10  Dice = 10
	D12  Dice = 12
	D20  Dice = 20
	D100 Dice = 100
)

// Roll rolls the given dice and returns the result.
func (d Dice) Roll() int {
	return 1 + rand.Intn(int(d))
}

// Roll rolls a number of dice and returns the sum of the results (plus the modifier).
func Roll(modifier int, dd Dice...) int {
	result := modifier
	for _, d := range dd {
		result += d.Roll()
	}
	return result
}

// RollAdvantage rolls a dice twice and returns the highest result (plus the modifier).
func RollAdvantage(modifier int, d Dice) int {
	return max(d.Roll(), d.Roll()) + modifier
}

// RollDisadvantage rolls a dice twice and returns the lowest result (plus the modifier).
func RollDisadvantage(modifier int, d Dice) int {
	return min(d.Roll(), d.Roll()) + modifier
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
