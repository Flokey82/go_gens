# genheightmap

This package provides various convenience functions to generate heightmaps.
The functions in this package are used for example by genmapvoronoi for the heightmap generation.

Generation of landscape features:
* Opensimplex noise
* Slope
* Cone
* Volcano cone
* Mountains/Hills

Operations on heightmaps:
* Normalization
* Relaxing
* Peakify (agitation / roughness)

## TODO

* Tidy up the code
* User-defined offsets
* User-defined seeds (or custom rand source)