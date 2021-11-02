# Simple Village Generator
This package implements the simplest approach for generating a somewhat stable and self-sustaining settlement.

## Principle

You instantiate a new BuildingPool, which will provide all available building types and add a number of types to it.

Then, you create a new settlement and give it a number of buildings you'd like to seed the economy with. After calling Solve(), the settlement should have a number of additional buildings needed to sustain itself... If there are still missing resources that the known building types cannot provide, these will have to be imported from another settlement.

## Example

```golang
package main

import (
	"github.com/Flokey82/go_gens/village"
	"log"
)

const (
	resFish   = "fish"
	resWorker = "worker"
	resGrain  = "grain"
	resFlour  = "flour"
	resBread  = "bread"
)

func main() {
	// Create a new building pool and register the known
	// building types.
	p := village.NewBuildingPool()
	bFishery := village.NewBuildingType("fishery")
	bFishery.Requires[resWorker] = 2
	bFishery.Provides[resFish] = 10
	p.AddType(bFishery)

	bHousing := village.NewBuildingType("housing")
	bHousing.Requires[resBread] = 4
	bHousing.Provides[resWorker] = 4
	p.AddType(bHousing)

	bFarm := village.NewBuildingType("farm")
	bFarm.Requires[resWorker] = 1
	bFarm.Provides[resGrain] = 10
	p.AddType(bFarm)

	bMill := village.NewBuildingType("mill")
	bMill.Requires[resGrain] = 10
	bMill.Requires[resWorker] = 1
	bMill.Provides[resFlour] = 10
	p.AddType(bMill)

	bBakery := village.NewBuildingType("bakery")
	bBakery.Requires[resFlour] = 2
	bBakery.Requires[resWorker] = 1
	bBakery.Provides[resBread] = 8
	p.AddType(bBakery)

	// Create a new settlement and add the fishery to seed
	// the economy.
	v := village.NewSettlement(p)
	v.AddBuilding(bFishery.NewBuilding())

	// Resolve all resource dependencies if possible.
	v.Solve()

	// Print all buildings of the settlement.
	log.Println(v.Buildings)
}
```
