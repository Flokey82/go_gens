# genmapvoxel

This package implements a simple voxel based terrain generator. Like... really simple.

Please be aware that this is just code for experimentation and is pretty rubbish for now... It will also serve as input for the genmarchingcubes package, so I can test out how well the logic for interpolation works and whatnot.

Honestly, I don't really feel motivated to write more right now.

NOTE: The stupid idea to have the origin be the center of the cube at 0,0 makes it more cumbersome to shrink the cubes vertically, since the Z axis for a standard cube ranges from -0.5 to +0.5. Yikes. I might change that.

## Done

* binary / bool voxel terrain
* float values for voxel data (-ish)
* hide faces that aren't visible
* wavefront OBJ export

## TODO

* merge faces of connected voxels

Here, have a picture:

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmapvoxel/images/blocky.png "Example of rendered voxel terrain (blocky)")


![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmapvoxel/images/smooth.png "Example of rendered voxel terrain (smooth)")