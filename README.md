# Go Generators

This repository contains various small attempts at procedural generation, simulation, and other things that might be useful for some people.

PLEASE NOTE: These packages are under constant development for maybe a month or two until I've figured out a consistent style etc.

DO NOT USE YET!

## Simulation and procedural generation (WIP)

### gamecs: Simulation with agents

This is a very basic simulation with agents using state machines and behavior trees.

### gamehex

Sample project to demonstrate how to use the Ebiten game engine to create a hexagonal game. This is still a work in progress.

![alt text](/gamehex/images/rgb.png "Hexagonal Game Demo")

### gameisometric

This is a simple isometric 'game' experiment using the ebitengine and is written in Go.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gameisometric/images/screenshot.png "Screenshot of isometric stuff!")

### gameraycast: Raycast implementation

This is a very basic raycast implementation using a randomly generated dungeon for testing. This isn't very feature rich yet, but it's a start.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gameraycast/images/textured.png "Screenshot of textured raycast!")

### gamerogueish: Rogue-like game skeleton

This is a very basic rogue-like game skeleton... Barely functional.

![alt text](https://raw.githubusercontent.com/Flokey82/gamerogueish/master/images/rgb.png "rogue-ish")

Now found here: [github.com/Flokey82/gamerogueish](https://gigithub.com/Flokey82/gamerogueish)

### genarchitecture: Architecture generator

This is a package for architectural style generation. It is still in the experimental phase and will probably change a lot in the future.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genarchitecture/images/rules.png "Generated mesh!")

### genbiome: Biome helper functions

This package provides functions for looking up the biome (e.g. forest, grassland, etc.) for the given precipitation and average temperature.

### gencellular: Simple cellular automata in Golang

This package currently only implements Conway's Game of Life.

### gencitymap: City map generator

This package is based on various implementations of city map generation. Currently supports simple random networks and tensor field based road network generation.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gencitymap/images/basic.png "Screenshot of first map!")

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gencitymap/images/tensor.png "Screenshot of basic tensor field!")

### gendemographics: Medieval demographics generator

This package is based on Medieval Demographics Made Easy by S. John Ross and generates demographics based on settlement density, population, etc.
The code is currently very messy and needs to be cleaned up.

### gendungeon: Dungeon generator

This package is a modified fork of https://github.com/brad811/go-dungeon/, which is an implementation of http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/gendungeon/images/lvl0.png "Multilevel Dungeon!")

This does now support multi-level dungeon generation with stairs.

### genempire: Empire generator (WIP)

This is a simple generator for empire names (based on genreligion).

### genfloortxt: Simple 2D text floorplan renderer

This package provides a little demo on how to 'render' walls as unicode characters, imitating the cp437 symbol set.

### genfurnishing: Simple generator for furnishing of rooms

This package will provide logic for generating clutter and furnishing for rooms of different types.

### genheightmap: Heightmap generator helper functions

This package provides some helper functions for generating heightmaps like slopes, cones, peakify, relax, etc.

### genlanguage: Language generator

Fantasy language generation wrapper used for fantasy map naming stuff.

### genlsystem: L-system generator

This package implements a 2d and 3d L-system generator with turtle graphics. This code is heavily based on:
* https://github.com/der-antikeks/lindenturtle/
* https://github.com/yalue/l_system_3d

### genmap2d: Simple 2D map generator in Golang (WIP)

This package provides functionality to generate 2D maps using procedural methods. It has a rudimentary settlement placement logic and probably needs lots work.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmap2d/images/rgb.png "Map!")

### genmap2derosion: Heightmap generation and erosion

This package is based on the fantastic work of Nick McDonald!
See: https://github.com/weigert/SimpleHydrology

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmap2derosion/images/screenshot.png "Eroded map")

### genmapvoronoi: Voronoi map generator

This is based on https://mewo2.com/notes/terrain/ and partially a port of https://github.com/mewo2/terrain/ to Go.

image: ![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmapvoronoi/images/obj_export.png "Map!")


### genmapvoxel: Voxel map generator

This package implements a simple voxel based terrain generator. Like... really simple. Don't get excited.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmapvoxel/images/smooth.png "Example of rendered voxel terrain (smooth)")

### genmarchingcubes: Marching cubes (WIP)

This package is an implementation of the marching cubes algorithm, which takes a voxel field and generates a surface mesh. Right now it is an almost straight fork of https://github.com/soypat/sdf/blob/main/render/marchingcubes.go and https://github.com/fogleman/mc/blob/master/mc.go but will be heavily modified in future versions. (It might already be, depending how lazy I am updating the READMEs)

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmarchingcubes/images/marched_float.png "Example of marched voxel terrain using float for fractional voxels")

### genmarchingsquares: Marching squares (WIP)

This package provides a simple implementation of the marching squares algorithm for illustrative purposes for now.

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/genmarchingsquares/images/squares.png "Example of rendered tiles")

### genreligion: Religion generator

This is based on Azgaars Fantasy Map Generator.

### genstory: Flavor text generator

Provides some functionality to create configurations for generating random flavor text.

### genthesaurus: Word association index

A very basic word association index... Work in progress.

### genworldvoronoi: Graph based planetary map generator

It simulates (somewhat) global winds and attempts to calculate precipitation and temperature for more intricate simulations in the future.
It features SVG, PNG, and Wavefront OBJ output.
This is based on https://www.redblobgames.com/x/1843-planet-generation/ and a port of https://github.com/redblobgames/1843-planet-generation to Go. 

![alt text](https://raw.githubusercontent.com/Flokey82/genworldvoronoi/master/images/relief_2.png "Relief Maps!")

### genvillage: Village economy generator

This package attempts to generate a stable, self sustaining economy given some initial buildings that serve as seed. The user can define building types that consume and produce resources, which the simulation will use to add new buildings until the sum of produced resources outweighs the sum of consumed resources. Please note that the output of this process is only meant to be used as input for a more sophisticated generator, only ensuring that besides required import of unavailable resources, the village/settlement can be in theory self-sustaining.

Don't get your expectations up, it's really basic. : P

### simmarket: Market economy simulation

Based on a few, way better implementations for economy simulation.

### simnpcs: NPCs with daily routines

This is an early experiment in simulating NPCs with daily routines. Sadly it is quite messy and doesn't follow best practices. It also has plenty of leftover code that I don't want to delete. I might clean it up in the future.

### simnpcs2: NPCs with daily routines (WIP)

This is a rewrite of simnpcs and gamecs... still in early stages of development.

### simvillage: Village simulation

This package is a port of the wonderful village simulator by Kontari. 
See: https://github.com/Kontari/Village/
Since I ported it straight from Python to Go, it is quite chaotic and probably would benefit from a refactor.

### simvillage_simple: Village simulation (simple)

A very basic village simulator which has settlers settle, procreate, live, love, and perish :) It is a very dumbed down version of kontari's village simulator, written from scratch.

### simvillage_tiles: Village simulation on a tile based map (WIP)

Currently, this only provides a very basic tile based renderer with some basic collision detection, etc. In future, it will be a village simulation on a tile based map. (Uses https://ebiten.org/)

![alt text](https://raw.githubusercontent.com/Flokey82/go_gens/master/simvillage_tiles/images/rgb.png "Screenshot!")

## Utility

### utils

This package contains a number of useful and commonly used functions and structures.

### vectors

Various vector math functions.

### vmesh

A voronoi mesh generation helper library. This is still very much bound to the genmapvoronoi package, so not very helpful for generic use. Sorry about that.
