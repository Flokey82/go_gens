package gamerogueish

import "github.com/BigJk/ramen/console"

type Scene interface {
	console.Component
	Close() error
}
