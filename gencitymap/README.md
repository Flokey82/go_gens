# gencitymap

This is a small tool to generate a city map based on simple rules wrt road branching and such.
Right now, it is just a proof of concept, but in the future it'll offer custom configurations, different map styles, and so on.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gencitymap/images/basic.png "Screenshot of first map!")

## Tensor Fields

This is based on the work of these folks:
https://github.com/ProbableTrain/MapGenerator

NOTE: HEAVILY WIP!! Do not use this yet!

### TODO

- [ ] Add coastline, water, rivers, etc.
- [X] Add graph generation from streamlines
- [X] Add polygon extraction for identifying plots and buildings
- [ ] Make code less buggy (crashes constantly)
- [ ] Add more map styles
- [ ] Eliminate artifacts in polygon extraction
- [ ] Investigate polygon extraction failure in polygons that contain dangling roads

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gencitymap/images/tensor.png "Screenshot of basic tensor field!")