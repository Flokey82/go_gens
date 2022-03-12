package gendemographics

// GenBusinesses returns a map of business to number of businesses for the given population size.
func GenBusinesses(population int) map[string]int {
	res := make(map[string]int)
	for _, bt := range BusinessTypes {
		if population > bt.Serves {
			//log.Println(fmt.Sprintf("%q: %d", bt.Name, population/bt.Serves))
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
	Name   string
	Serves int
}

// New returns a new business of the given business types.
func (bt *BusinessType) New() *Business {
	return &Business{
		BusinessType: bt,
	}
}

var BusinessTypes = []*BusinessType{
	{"Shoemaker", 150},
	{"Furrier", 250},
	{"Maidservant", 250},
	{"Tailor", 250},
	{"Barber", 350},
	{"Jeweler", 400},
	{"Taverns/Restaurant", 400},
	{"Old-clothes", 400},
	{"Pastrycook", 500},
	{"Mason", 500},
	{"Carpenter", 550},
	{"Weaver", 600},
	{"Chandler", 700},
	{"Mercer", 700},
	{"Cooper", 700},
	{"Baker", 800},
	{"Watercarrier", 850},
	{"Scabbardmaker", 850},
	{"Wine-seller", 900},
	{"Hatmaker", 950},
	{"Saddler", 1000},
	{"Chicken Butcher", 1000},
	{"Pursemaker", 1100},
	{"Woodseller", 2400},
	{"Magic-shop", 2800},
	{"Bookbinder", 3000},
	{"Butcher", 1200},
	{"Fishmonger", 1200},
	{"Beer-seller", 1400},
	{"Buckle Maker", 1400},
	{"Plasterer", 1400},
	{"Spice Merchant", 1400},
	{"Blacksmith", 1500},
	{"Painter", 1500},
	{"Doctor", 1700},
	{"Roofer", 1800},
	{"Locksmith", 1900},
	{"Bather", 1900},
	{"Ropemaker", 1900},
	{"Inn", 2000},
	{"Tanner", 2000},
	{"Copyist", 2000},
	{"sculptor", 2000},
	{"Rugmaker", 2000},
	{"Harness-Maker", 2000},
	{"Bleacher", 2100},
	{"Hay Merchant", 2300},
	{"Cutler", 2300},
	{"Glovemaker", 2400},
	{"Woodcarver", 2400},
	{"Bookseller", 6300},
	{"Illuminator", 3900},
	{"Place Of Worship", 400},
	// OTHER JOBS / NOT BUSINESS
	{"Law Enforcement", 150},
	{"Noble", 200},
	{"Administrator", 650},
	{"Clergy", 40},
	{"Priest", 1000},
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
