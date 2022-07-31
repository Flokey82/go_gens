# Go Generators
This repository contains various small attempts at procedural generation, simulation, and other things that might be useful for some people.

PLEASE NOTE: These packages are under constant development for maybe a month or two until I've figured out a consistent style etc.

DO NOT USE YET!

## Simulation and procedural generation (WIP)

### gencellular: Simple cellular automata in Golang
This package currently only implements Conway's Game of Life.

### gendemographics: Medieval demographics generator
This package is based on Medieval Demographics Made Easy by S. John Ross.

### gendungeon: Dungeon generator
This package is a modified fork of https://github.com/brad811/go-dungeon/, which is an implementation of http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/.

### genvillage: Village economy generator
This package attempts to generate a stable, self sustaining economy given some initial buildings that serve as seed. The user can define building types that consume and produce resources, which the simulation will use to add new buildings until the sum of produced resources outweighs the sum of consumed resources. Please note that the output of this process is only meant to be used as input for a more sophisticated generator, only ensuring that besides required import of unavailable resources, the village/settlement can be in theory self-sustaining.

### genfloortxt: Simple 2D text floorplan renderer
This package provides a little demo on how to 'render' walls as unicode characters, imitating the cp437 symbol set.

### genmap2derosion: Heightmap generation and erosion
This package is based on the fantastic work of Nick McDonald!
See: https://github.com/weigert/SimpleHydrology

### genmapvoronoi: Voronoi map generator
This is based on https://mewo2.com/notes/terrain/ and partially a port of https://github.com/mewo2/terrain/ to Go.

### genworldvoronoi: Graph based planetary map generator
This is based on https://www.redblobgames.com/x/1843-planet-generation/ and a port of https://github.com/redblobgames/1843-planet-generation to Go.

### simmotive: Sims motive
This package is a crude port of https://github.com/alexcu/motive-simulator which is adapted from Don Hopkins' article The Soul of The Sims which shows an prototype of the 'soul' of what became The Sims 1, written January 23, 1997.

### simvillage: Village simulation
This package is a port of the wonderful village simulator by Kontari. See: https://github.com/Kontari/Village/

### simvillage_simple: Village simulation (simple)
A very basic village simulator which has settlers settle, procreate, live, love, and perish :)

## Game components

### gameloop: Simple Game Loop
This is a very, very basic game loop... nothing fancy about it.