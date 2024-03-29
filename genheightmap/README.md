# genheightmap

This package provides various convenience functions to generate heightmaps.
The functions in this package are used for example by genmapvoronoi for the heightmap generation.

Generation of landscape features:
* Opensimplex noise
* Slope
* Cone
* Volcano cone
* Crater
* Fissure
* Mountains/Hills

Operations on heightmaps:
* Normalization
* Relaxing
* Peakify (agitation / roughness)

![alt text](/genheightmap/images/crater.png "crater")

![alt text](/genheightmap/images/fissure.png "fissure")

![alt text](/genheightmap/images/mountainrange.png "mountain range")

![alt text](/genheightmap/images/volcone.png "volcone")

![alt text](/genheightmap/images/hills.png "hills")

## TODO

* Tidy up the code
* User-defined offsets
* User-defined seeds (or custom rand source)