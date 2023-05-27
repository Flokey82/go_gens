# Hydraulic Erosion
This package is based on the fantastic work of Nick McDonald!
See: https://github.com/weigert/SimpleHydrology

## Principle

See: https://nickmcd.me/2020/04/15/procedural-hydrology/
See: https://nickmcd.me/2020/11/23/particle-based-wind-erosion/

## Notes

This is not a complete port of the code mentioned above and includes some experimental alternatives for determining water flux information.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmap2derosion/images/screenshot.png "Eroded map")

## TODO
* Fix up climate simulation
  * Use primary heightmap directly
* Finalize flood algorithm documentation
* Either complete or remove flux based hydrology
* Add vegetation
* Clean up different erosion algorithms
* Make effective use of sediment in hydrology

## Done
* Get rid of terrain struct
* Move biomes to climate struct
* Basic wind erosion
