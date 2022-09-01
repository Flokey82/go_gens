# genworldvoronoi: Graph based planetary map generator

It simulates (somewhat) global winds and attempts to calculate precipitation and temperature for more intricate simulations in the future.
It features SVG, PNG, and Wavefront OBJ output.

This is based on https://www.redblobgames.com/x/1843-planet-generation/ and a port of https://github.com/redblobgames/1843-planet-generation to Go. 

I draw further inspiration from various other projects such as https://github.com/weigert/SimpleHydrology and https://github.com/mewo2/terrain

... if you haven't noticed yet, this is a placeholder for when I feel less lazy and add more information :D

## TODO

* Re-evaluate rainfall and moisture distribution
* Climate is surprisingly humid along equator?

Here some old pictures what it does...

## SVG export with rivers, capital city placement and stuff.
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/svg.png "Screenshot of SVG!")

## Does political maps.
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/political.png "Political Maps!")

## Simulates climate (-ish)
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/climate.png "Screenshot of Biomes!")

## Exports to Wavefront OBJ
![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genworldvoronoi/images/obj.png "Screenshot of OBJ Export in Blender!")