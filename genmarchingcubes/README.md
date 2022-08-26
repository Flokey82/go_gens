# genmarchingcubes

This package is an implementation of the marching cubes algorithm, which takes a voxel field and generates a surface mesh.

Right now it is almost a straight fork of
* https://github.com/soypat/sdf/blob/main/render/marchingcubes.go
* https://github.com/fogleman/mc/blob/master/mc.go

... but will be modified more heavily in future versions.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmarchingcubes/images/voxel.png "Example of source voxel terrain")

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmarchingcubes/images/marched.png "Example of marched voxel terrain")

## References
* https://en.wikipedia.org/wiki/Marching_cubes