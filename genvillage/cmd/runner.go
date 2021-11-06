package main

import (
	"github.com/Flokey82/go_gens/genvillage"
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
	p := genvillage.NewBuildingPool()
	bFishery := genvillage.NewBuildingType("fishery")
	bFishery.Requires[resWorker] = 2
	bFishery.Provides[resFish] = 10
	p.AddType(bFishery)

	bHousing := genvillage.NewBuildingType("housing")
	bHousing.Requires[resBread] = 4
	bHousing.Provides[resWorker] = 4
	p.AddType(bHousing)

	bFarm := genvillage.NewBuildingType("farm")
	bFarm.Requires[resWorker] = 1
	bFarm.Provides[resGrain] = 10
	// FIXME: A farm usually doubles as housing, but we don't check
	// if buildings provide a net-positive on provided resources yet, which
	// might lead to many farms being randomly added when workers are needed.
	// bFarm.Provides[resWorker] = 1
	p.AddType(bFarm)

	bMill := genvillage.NewBuildingType("mill")
	bMill.Requires[resGrain] = 10
	bMill.Requires[resWorker] = 1
	bMill.Provides[resFlour] = 10
	p.AddType(bMill)

	bBakery := genvillage.NewBuildingType("bakery")
	bBakery.Requires[resFlour] = 2
	bBakery.Requires[resWorker] = 1
	bBakery.Provides[resBread] = 8
	p.AddType(bBakery)

	// Create a new settlement and add the fishery to seed
	// the economy.
	v := genvillage.NewSettlement(p)
	v.AddBuilding(bFishery.NewBuilding())

	// Resolve all resource dependencies if possible.
	v.Solve()

	// Print all buildings of the settlement.
	log.Println(v.Buildings)
}
