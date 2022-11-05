# genbiome: Biome helper functions

This package provides functions for looking up the biome (e.g. forest, grassland, etc.) for the given precipitation and average temperature.

There are two (and a half) variants:

## Redblobgames' biome map

This biome map is based on http://www-cs-students.stanford.edu/~amitp/game-programming/polygon-map-generation/#biomes and rather than temperature and precipitation it uses elevation- and moisture zones.
This makes it a bit simpler, but also coarse.

## Whittaker's biome map

This map is based on https://github.com/JoeyR/FastBiome and provides a more extensive and fine(r) grained biome map. It is modeled after real-world precipitation and temperature corellations as described by Whittaker (see: https://en.wikipedia.org/wiki/Biome#Whittaker_(1962,_1970,_1975)_biome-types).

### Extended variant

I've started to extend the Whittaker model to be more complete and have added the missing biomes, as well as additional biomes such as "snow", "wetlands", and "hot swamp". I am not sure if I will keep the swamp though, since it is a bit... meh. Let's say I don't really trust my skills as a climate biome science person.

## Azgaar's biome map

This is a third variant, which is based on Azgaar's map generator and its biome map. It is a bit more fine grained than the Redblobgames' biome map, but also a bit more coarse than the Whittaker map. It is also based on elevation and moisture zones, but has a few more biomes, and takes temperature into account. (see: https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/main.js)

## Done

* missing biome in extended Whittaker (cool, humid)

## TODO

* improve function names
* min/max values should be ranges
* find better colors for biomes
  * wetlands, savannah, ... everything?
  * create a color palette?
* re-do extended Whittaker (seems to be off)
  * wetlands covers a too broad range
  * swamps need to be looked at