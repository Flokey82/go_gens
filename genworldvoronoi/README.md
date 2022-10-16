# genworldvoronoi: Graph based planetary map generator

It simulates (somewhat) global winds and attempts to calculate precipitation and temperature for more intricate simulations in the future.
It features SVG, PNG, and Wavefront OBJ output.

This is based on https://www.redblobgames.com/x/1843-planet-generation/ and a port of https://github.com/redblobgames/1843-planet-generation to Go. 

I draw further inspiration from various other projects such as https://github.com/weigert/SimpleHydrology and https://github.com/mewo2/terrain

... if you haven't noticed yet, this is a placeholder for when I feel less lazy and add more information :D

## TODO

* Climate is surprisingly humid along equator?
  * Note: That is actually correct... the dryest climate is around +30° latitude, which is the subtropical desert. The equator is not dry at all.
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
* Resources
  * Improve resource distribution
    * Fitness functions
    * Amount / quality / discoverability?
  * Add more resource types

Here some old pictures what it does...

## SVG export with rivers, capital city placement and stuff.
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/svg.png "Screenshot of SVG!")

## Does political maps.
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/political.png "Political Maps!")

## Simulates climate (-ish)
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/climate.png "Screenshot of Biomes!")

## Exports to Wavefront OBJ
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/obj.png "Screenshot of OBJ Export in Blender!")