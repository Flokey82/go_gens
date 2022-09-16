package gameconstants

// Grain yield of kg grain per m^2 per year (modern).
// TODO: Verify with second source.
// https://www.regenerativegardener.org/442014395
// WheatEfficiency  = 0.5   // kg/m^2
// OatEfficiency    = 0.3   // kg/m^2
// BarleyEfficiency = 0.4   // kg/m^2
// CornEfficiency   = 0.05  // kg/m^2
// RyeEfficiency    = 0.017 // kg/m^2
// QuinoaEfficiency = 0.5   // kg/m^2
// MilletEfficiency = 0.2   // kg/m^2

// Crop yields in the middle ages:
// http://piketty.pse.ens.fr/files/Allen2005.pdf
//
// Gross yields, bushels per acre
//
//             1300  1500  1700  1750  1800  1850
// ----------------------------------------------
// Wheat       10.7  14.0  19.0  20.0  20.0  28.0
// Rye         10.4  14.0  19.0  20.0  22.0  28.0
// Barley      16.3  17.0  29.0  30.0  30.0  35.0
// Oats        11.3  14.0  21.0  35.0  38.0  40.0
// Beans/peas  11.8  14.0  19.0  20.0  20.0  28.0
// Potatoes       –     – 150.0 150.0 150.0 150.0

// Conversion bushels to kg
// https://www.sagis.org.za/conversion_table.html
// http://www.solluscapital.com.br/English/conversion_table.html
//
// Wheat & Soybeans: 1 bushel = 27.216 kg
// Barley:           1 bushel = 21.772 kg
// Maize & Sorghum:  1 bushel = 25.401 kg
// Oats:             1 bushel = 14.515 kg
// Corn:             1 bushel = 25.401 kg

// https://www.smallfarmcanada.ca/resources/standard-weights-per-bushel-for-agricultural-commodities
// Extensive list above ...
//
// Potatoes:         1 bushel = 27.2155 kg
// Peas:			 1 bushel = 14.5150 kg
// Rye:			     1 bushel = 25.4012 kg

// https://www.jstor.org/stable/2595059 (corn yields)
// https://www.basvanleeuwen.net/bestanden/agriclongrun1250to1850.pdf
// https://github.com/jgardner1/pysgame1/tree/addf619f87fc51850eca504bc2271fc3bcbea02a/docs

// Water requirement etc
// https://github.com/AmericasWater/awash/blob/810ce8ad229e54b5cecfc6ffc45c5b48040a86f0/src/lib/agriculture.jl
// https://github.com/jdischler/GrazeScapeServer/tree/bcafdc2f7b22cedcc72fb42ca7d20b546e9c958d/conf/modelDefs/yield

// Medieval agriculture
// https://worldbuilding.stackexchange.com/questions/125445/size-of-family-owned-medieval-farm
// https://www.ibiblio.org/london/agriculture/general/1/msg00070.html

// Sustainable Agriculture in the Middle Ages: The English Manor
// https://www.jstor.org/stable/40274710
// https://eprints.bbk.ac.uk/id/eprint/17662/1/sustainable.pdf
//
// Produce from a typical manorial estate
//
// Source                     | Produce
// ====================================================================
// Arable Crops
// --------------------------------------------------------------------
// Wheat (spelt, club, bread) | Bread, ale
// Oats (cultivated and wild) | Bread, pottage, livestock feed, ale
// Barley (hulled, naked)     | Ale, bread, livestock feed
// Rye                        | Bread
// Peas, beans, vetches       | Whole plant for human and livestock food
// All cereal straw           | Livestock feed, thatching
//
// Orchard and Garden Crops
// --------------------------------------------------------------------
// Apples                     | Fruit, cider
// Pears, cherries, figs,     | Fruit and nuts
// walnuts, damsons, plums    |
// Vines                      | Wine
// Flax                       | Linen
// Hemp                       | Rope and linen
// Herbs                      | Seasoning, medicines, dyes
// Leeks, onions, borage,     | Vegetable foods
// mustard, peas, beans       |
//
// Livestock
// --------------------------------------------------------------------
// Pigs                       | Meat
// Sheep, goats               | Wool, milk, manures,
//                            | some meat, skin for
//                            | parchment
// Cattle                     | Draught power, milk,
//                            | cheese, butter, curds,
//                            | some meat, leather, horn
// Horses                     | Draught power, leather
// Poultry (chickens, geese,  | Eggs, meat
// swans, peacocks)           |
// Pigeons and doves          | Meat, manures
// Bees                       | Honey, wax
// Rabbits                    | Meat, fur
//
// Natural Resources
// --------------------------------------------------------------------
// Deer                       | Meat, manures
// Wild boar                  | Meat
// Birds                      | Meat
// Fish - from fish pond,     | Meat
// river, sea                 |
// Hares                      | Meat, fur
// Oak and beech trees        | Acorns and mast for pigs,
//                            | timber
// Other trees and shrubs     | Nuts, berries, fruits,
//                            | timber, browse, fuelwood
// Ferns, bracken, sedges     | Thatch, bedding, litter
// Nettles                    | Linen
// Osiers, reeds              | Baskets, fish traps
// Holly, thorns              | Threshing flails
// Peat                       | Fuel
// Herbs                      | Medicines, vegetables
// Grass                      | Hay
// Grass turves               | Roofing, fuel
//
// * Lord Ernle, Enrllish Farthing. Past and Present, 6th edn, 196I, pp 6-30. Gras and Gras, op cit, pp 33-53. H S Bennett, L~' on the English
// Manor. A Study of Peasant Conditions, 115o-14oo, Cambridge, 1937, pp 75-96. G W Johnson, A History of Gardening, 1829, pp 36-43.
// J Harvey, Medieval Gardens, 198 i, pp 163-I8o; E M Veale, 'The Rabbit in England', Ag Hist Rev, V, 1957, pp 85-90.

// Ploughing speed, harrowing speed.
// https://eprints.bbk.ac.uk/id/eprint/17662/1/sustainable.pdf
// PloughingSpeedOxenMin  = 0.1 // (ha/day)
// PloughingSpeedOxenMax  = 0.4 // (ha/day)
// PloughingSpeedHorseMin = 0.3 // (ha/day)
// PloughingSpeedHorseMax = 0.5 // (ha/day)
//
// Oat consumption.
// OatConsumptionOats   = 61.0  // (kg/year)
// HorseConsumptionOats = 362.0 // (kg/year) Yes, that seems to be correct.

// The Economics of Horses and Oxen in Medieval England
// https://bahs.org.uk/AGHR/ARTICLES/30n1a3.pdf
// https://www.jstor.org/stable/40274182

// THE PRODUCTIVITY OF PEASANT AGRICULTURE
// Average Oakington demesne yields, net of seed, in bushels, 1361–99
// Wheat Maslin Dredge White peas Black peas Oats
// https://www.jstor.org/stable/42921567
