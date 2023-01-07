# genworldvoronoi: Graph based planetary map generator

It simulates (somewhat) global winds and attempts to calculate precipitation and temperature for more intricate simulations in the future.
It features SVG, PNG, and Wavefront OBJ output.

This is based on https://www.redblobgames.com/x/1843-planet-generation/ and a port of https://github.com/redblobgames/1843-planet-generation to Go. 

I draw further inspiration from various other projects such as https://github.com/weigert/SimpleHydrology and https://github.com/mewo2/terrain

... if you haven't noticed yet, this is a placeholder for when I feel less lazy and add more information :D

## TODO

* Climate is surprisingly humid along equator?
  * Note: That is actually correct... the dryest climate is around +30° latitude, which is the subtropical desert. The equator is not dry at all.
  * Add desert oases that are fed from underground aquifers. Look at these examples: https://www.google.com/maps/d/viewer?mid=1BvY10l3yzWt48IwCXqDcyeuawpA&hl=en&ll=26.715853962142784%2C28.408963168787885&z=6
  * Climate seems too wet (too many wetlands?)
  * Seasonal forests should not be at the equator, where there are no seasons.
* Elevation
  * The heightmap is currently generated as a linear interpolation of three distance fields. This results in relative unrealistic elevation distribution.
* Winds
  * Make winds push temperature around, not just humidity
  * Re-evaluate rainfall and moisture distribution
  * Check if we also push dry air around, not just humid air
* Civilization
  * Industry and trade
    * Improve trade routes
    * Introduce industry
  * Cities
    * Add better city fitness functions
    * Separate world generation better from everything else
    * Assign goods and resources to cities (for trade)
  * Empires
    * Provide simpler means to query information over an empire
    * Introduce empires with capitals
  * Cultures
    * Add fitness function for "natural" population density estimates
* Resources
  * Improve resource distribution
    * Fitness functions
    * Amount / quality / discoverability?
  * Add more resource types

Here some old pictures what it does...

## SVG export with rivers, capital city placement and stuff.
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/svg.png "Screenshot of SVG!")

## Leaflet server (and sad flavor text).
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/leaflet.png "Flavortext Maps!")

## Poor man's relief map.
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/relief.png "Relief Maps!")

## Does political maps.
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/political.png "Political Maps!")

## Simulates climate (-ish)
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/climate.png "Screenshot of Biomes!")

## Simulates seasons (-ish)
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/seasons.webp "Screenshot of Seasons!")

## Exports to Wavefront OBJ
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/obj.png "Screenshot of OBJ Export in Blender!")