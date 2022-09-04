package gendemographics

// GenBusinesses returns a map of business to number of businesses for the given population size.
// NOTE: Meh, not happy about it being a map.
func GenBusinesses(population int) map[string]int {
	res := make(map[string]int)
	for _, bt := range BusinessTypes {
		if population > bt.Serves {
			res[bt.Name] = population / bt.Serves
		}
	}
	return res
}

// Business represents a single business of a given BusinessType.
type Business struct {
	*BusinessType
}

// BusinessType represents a type of business.
type BusinessType struct {
	Name   string // business type name (profession)
	Serves int    // people served by a single business
}

// New returns a new business of the given business types.
func (bt *BusinessType) New() *Business {
	return &Business{
		BusinessType: bt,
	}
}

// BusinessTypes is a list of all business types.
//
// NOTE: Wouldn't this be better in a CSV or something?
// On the other hand, it'd be good to have some form of
// identifier that users of the package can refer to.
// IDK.
//
// ... should it be a map?
var BusinessTypes = []*BusinessType{{
	Name:   "Shoemaker",
	Serves: 150,
}, {
	Name:   "Furrier",
	Serves: 250,
}, {
	Name:   "Maidservant",
	Serves: 250,
}, {
	Name:   "Tailor",
	Serves: 250,
}, {
	Name:   "Barber",
	Serves: 350,
}, {
	Name:   "Jeweler",
	Serves: 400,
}, {
	Name:   "Taverns/Restaurant",
	Serves: 400,
}, {
	Name:   "Old-clothes",
	Serves: 400,
}, {
	Name:   "Pastrycook",
	Serves: 500,
}, {
	Name:   "Mason",
	Serves: 500,
}, {
	Name:   "Carpenter",
	Serves: 550,
}, {
	Name:   "Weaver",
	Serves: 600,
}, {
	Name:   "Chandler",
	Serves: 700,
}, {
	Name:   "Mercer",
	Serves: 700,
}, {
	Name:   "Cooper",
	Serves: 700,
}, {
	Name:   "Baker",
	Serves: 800,
}, {
	Name:   "Watercarrier",
	Serves: 850,
}, {
	Name:   "Scabbardmaker",
	Serves: 850,
}, {
	Name:   "Wine-seller",
	Serves: 900,
}, {
	Name:   "Hatmaker",
	Serves: 950,
}, {
	Name:   "Saddler",
	Serves: 1000,
}, {
	Name:   "Chicken Butcher",
	Serves: 1000,
}, {
	Name:   "Pursemaker",
	Serves: 1100,
}, {
	Name:   "Woodseller",
	Serves: 2400,
}, {
	Name:   "Magic-shop",
	Serves: 2800,
}, {
	Name:   "Bookbinder",
	Serves: 3000,
}, {
	Name:   "Butcher",
	Serves: 1200,
}, {
	Name:   "Fishmonger",
	Serves: 1200,
}, {
	Name:   "Beer-seller",
	Serves: 1400,
}, {
	Name:   "Buckle Maker",
	Serves: 1400,
}, {
	Name:   "Plasterer",
	Serves: 1400,
}, {
	Name:   "Spice Merchant",
	Serves: 1400,
}, {
	Name:   "Blacksmith",
	Serves: 1500,
}, {
	Name:   "Painter",
	Serves: 1500,
}, {
	Name:   "Doctor",
	Serves: 1700,
}, {
	Name:   "Roofer",
	Serves: 1800,
}, {
	Name:   "Locksmith",
	Serves: 1900,
}, {
	Name:   "Bather",
	Serves: 1900,
}, {
	Name:   "Ropemaker",
	Serves: 1900,
}, {
	Name:   "Inn",
	Serves: 2000,
}, {
	Name:   "Tanner",
	Serves: 2000,
}, {
	Name:   "Copyist",
	Serves: 2000,
}, {
	Name:   "sculptor",
	Serves: 2000,
}, {
	Name:   "Rugmaker",
	Serves: 2000,
}, {
	Name:   "Harness-Maker",
	Serves: 2000,
}, {
	Name:   "Bleacher",
	Serves: 2100,
}, {
	Name:   "Hay Merchant",
	Serves: 2300,
}, {
	Name:   "Cutler",
	Serves: 2300,
}, {
	Name:   "Glovemaker",
	Serves: 2400,
}, {
	Name:   "Woodcarver",
	Serves: 2400,
}, {
	Name:   "Bookseller",
	Serves: 6300,
}, {
	Name:   "Illuminator",
	Serves: 3900,
}, {
	Name:   "Place Of Worship",
	Serves: 400,
},
	// OTHER JOBS / NOT BUSINESS
	{
		Name:   "Law Enforcement",
		Serves: 150,
	}, {
		Name:   "Noble",
		Serves: 200,
	}, {
		Name:   "Administrator",
		Serves: 650,
	}, {
		Name:   "Clergy",
		Serves: 40,
	}, {
		Name:   "Priest",
		Serves: 1000,
	},
}

/* ORIGINAL LIST
 * Shoemakers 150
 * Furriers 250
 * Maidservants 250
 * Tailors 250
 * Barbers 350
 * Jewelers 400
 * Taverns/Restaurants 400
 * Old-clothes 400 ?
 * Pastrycooks 500
 * Masons 500
 * Carpenters 550
 * Weavers 600
 * chandlers 700
 * Mercers 700
 * Coopers 700
 * Bakers 800
 * Watercarriers 850 This is a job
 * Scabbardmakers 850
 * wine-sellers 900
 * Hatmakers 950
 * Saddlers 1,000
 * Chicken Butchers 1,000
 * Pursemakers 1,100
 * Woodsellers 2,400
 * Magic-shops 2,800
 * Bookbinders 3,000
 * Butchers 1,200
 * Fishmongers 1,200
 * Beer-sellers 1,400
 * Buckle Makers 1,400
 * Plasterers 1,400
 * Spice Merchants 1,400
 * Blacksmiths 1,500
 * Painters 1,500
 * Doctors* 1,700
 * Roofers 1,800
 * Locksmiths 1,900
 * Bathers 1,900
 * Ropemakers 1,900
 * Inns 2,000
 * Tanners 2,000
 * Copyists 2,000
 * sculptors 2,000
 * Rugmakers 2,000
 * Harness-Makers 2,000
 * Bleachers 2,100
 * Hay Merchants 2,300
 * Cutlers 2,300
 * Glovemakers 2,400
 * Woodcarvers 2,400
 * Booksellers 6,300
 * Illuminators 3,900
 * Place Of Worship 400
 *
 * OTHER JOBS / NOT BUSINESS
 *
 * Law Enforcement 150
 * Noble 200
 * Administrator 650
 * Clergy 40
 * Priest 1000
 *
 */
