# gendemographics

This package implements a simple demographic generator for games and is based on Medieval Demographics Made Easy by S. John Ross.

## Done

* Settlement demographics
  * Fixed settlement sizes (somewhat).
  * Added some variation

## TODO

* Business demographics
  * Make business types customizable
  * Find a better way to store the defaults (CSV?)
* Settlement demographics
  * Allow custom settlement types / sizes (town, village, hamlet, ...)
  * Allow custom farmland allotment per household
* General
  * Move constants to the same spot.
  * Move constants to config struct for customization?

## Usage

```go
package main

import (
	"fmt"

	"github.com/Flokey82/go_gens/gendemographics"
)

func main() {
	c := gendemographics.New()
	n := c.NewNation(41000, gendemographics.DensityMedium)
	n.Log()
}
```
