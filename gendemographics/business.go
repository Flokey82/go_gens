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

// NOTE: Wouldn't this be better in a CSV or something?
// On the other hand, it'd be good to have some form of
// identifier that users of the package can refer to.
// IDK.
//
// ... should it be a map?
var (
	BizShoeMaker      = &BusinessType{Name: "Shoemaker", Serves: 150}
	BizFurrier        = &BusinessType{Name: "Furrier", Serves: 250}
	BizMaidservant    = &BusinessType{Name: "Maidservant", Serves: 250}
	BizTailor         = &BusinessType{Name: "Tailor", Serves: 250}
	BizBarber         = &BusinessType{Name: "Barber", Serves: 350}
	BizJeweler        = &BusinessType{Name: "Jeweler", Serves: 400}
	BizTavern         = &BusinessType{Name: "Taverns/Restaurant", Serves: 400}
	BizOldClothes     = &BusinessType{Name: "Old-clothes", Serves: 400}
	BizPastryCook     = &BusinessType{Name: "Pastrycook", Serves: 500}
	BizMason          = &BusinessType{Name: "Mason", Serves: 500}
	BizCarpenter      = &BusinessType{Name: "Carpenter", Serves: 550}
	BizWeaver         = &BusinessType{Name: "Weaver", Serves: 600}
	BizChandler       = &BusinessType{Name: "Chandler", Serves: 700}
	BizMercer         = &BusinessType{Name: "Mercer", Serves: 700}
	BizCooper         = &BusinessType{Name: "Cooper", Serves: 700}
	BizBaker          = &BusinessType{Name: "Baker", Serves: 800}
	BizWaterCarrier   = &BusinessType{Name: "Water carrier", Serves: 850}
	BizScabbardmaker  = &BusinessType{Name: "Scabbardmaker", Serves: 850}
	BizWineSeller     = &BusinessType{Name: "Wine-seller", Serves: 900}
	BizHatmaker       = &BusinessType{Name: "Hatmaker", Serves: 950}
	BizSaddler        = &BusinessType{Name: "Saddler", Serves: 1000}
	BizChickenButcher = &BusinessType{Name: "Chicken Butcher", Serves: 1000}
	BizPursemaker     = &BusinessType{Name: "Pursemaker", Serves: 1100}
	BizWoodseller     = &BusinessType{Name: "Woodseller", Serves: 2400}
	BizMagicShop      = &BusinessType{Name: "Magic-shop", Serves: 2800}
	BizBookbinder     = &BusinessType{Name: "Bookbinder", Serves: 3000}
	BizButcher        = &BusinessType{Name: "Butcher", Serves: 1200}
	BizFishmonger     = &BusinessType{Name: "Fishmonger", Serves: 1200}
	BizBeerSeller     = &BusinessType{Name: "Beer-seller", Serves: 1400}
	BizBucklemaker    = &BusinessType{Name: "Bucklemaker", Serves: 1400}
	BizPlasterer      = &BusinessType{Name: "Plasterer", Serves: 1400}
	BizSpiceMerchant  = &BusinessType{Name: "Spice merchant", Serves: 1400}
	BizBlacksmith     = &BusinessType{Name: "Blacksmith", Serves: 1500}
	BizPainter        = &BusinessType{Name: "Painter", Serves: 1500}
	BizDoctor         = &BusinessType{Name: "Doctor", Serves: 1700}
	BizRoofer         = &BusinessType{Name: "Roofer", Serves: 1800}
	BizLocksmith      = &BusinessType{Name: "Locksmith", Serves: 1900}
	BizBather         = &BusinessType{Name: "Bather", Serves: 1900}
	BizRopemaker      = &BusinessType{Name: "Ropemaker", Serves: 1900}
	BizInn            = &BusinessType{Name: "Inn", Serves: 2000}
	BizTanner         = &BusinessType{Name: "Tanner", Serves: 2000}
	BizCopyist        = &BusinessType{Name: "Copyist", Serves: 2000}
	BizSculptor       = &BusinessType{Name: "Sculptor", Serves: 2000}
	BizRugmaker       = &BusinessType{Name: "Rugmaker", Serves: 2000}
	BizHarnessMaker   = &BusinessType{Name: "Harness-maker", Serves: 2000}
	BizBleacher       = &BusinessType{Name: "Bleacher", Serves: 2100}
	BizHayMerchant    = &BusinessType{Name: "Hay merchant", Serves: 2300}
	BizCutler         = &BusinessType{Name: "Cutler", Serves: 2300}
	BizGlovemaker     = &BusinessType{Name: "Glovemaker", Serves: 2400}
	BizWoodcarver     = &BusinessType{Name: "Woodcarver", Serves: 2400}
	BizBookseller     = &BusinessType{Name: "Bookseller", Serves: 6300}
	BizIlluminator    = &BusinessType{Name: "Illuminator", Serves: 3900}
	BizPlaceOfWorship = &BusinessType{Name: "Place of Worship", Serves: 400}
	BizLawEnforcement = &BusinessType{Name: "Law Enforcement", Serves: 150}
	BizNoble          = &BusinessType{Name: "Noble", Serves: 200}
	BizAdministrator  = &BusinessType{Name: "Administrator", Serves: 650}
	BizClergy         = &BusinessType{Name: "Clergy", Serves: 40}
	BizPriest         = &BusinessType{Name: "Priest", Serves: 1000}
)

// BusinessTypes is a list of all business types.
var BusinessTypes = []*BusinessType{
	BizShoeMaker,
	BizFurrier,
	BizMaidservant,
	BizTailor,
	BizBarber,
	BizJeweler,
	BizTavern,
	BizOldClothes,
	BizPastryCook,
	BizMason,
	BizCarpenter,
	BizWeaver,
	BizChandler,
	BizMercer,
	BizCooper,
	BizBaker,
	BizWaterCarrier,
	BizScabbardmaker,
	BizWineSeller,
	BizHatmaker,
	BizSaddler,
	BizChickenButcher,
	BizPursemaker,
	BizWoodseller,
	BizMagicShop,
	BizBookbinder,
	BizButcher,
	BizFishmonger,
	BizBeerSeller,
	BizBucklemaker,
	BizPlasterer,
	BizSpiceMerchant,
	BizBlacksmith,
	BizPainter,
	BizDoctor,
	BizRoofer,
	BizLocksmith,
	BizBather,
	BizRopemaker,
	BizInn,
	BizTanner,
	BizCopyist,
	BizSculptor,
	BizRugmaker,
	BizHarnessMaker,
	BizBleacher,
	BizHayMerchant,
	BizCutler,
	BizGlovemaker,
	BizWoodcarver,
	BizBookseller,
	BizIlluminator,
	BizPlaceOfWorship,
	BizLawEnforcement,
	BizNoble,
	BizAdministrator,
	BizClergy,
	BizPriest,
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
